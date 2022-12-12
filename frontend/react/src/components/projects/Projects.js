import {useContext, useEffect, useState} from "react";
import useFetch from "../../hooks/useFetch";
import {AuthContext} from "../../context/AuthContext";
import Loader from "../../uikit/Loader";
import {Routes, Route} from "react-router-dom";
import Project from "./Project";
import CreateProject from "./CreateProject";
import {AccessGuard, useAccess} from "@a11n-io/cerberus-reactjs"
import {Col, Container, ListGroup, ListGroupItem, Row} from "react-bootstrap";
import {ProjectContext} from "./ProjectContext";

export default function Projects() {

    return <>

        <Routes>
            <Route path="*" element={<ProjectList/>}/>
        </Routes>

    </>

}

function ProjectList() {
    const [projects, setProjects] = useState([])
    const [selectedProject, setSelectedProject] = useState(null)
    const authCtx = useContext(AuthContext)
    const projectCtx = useContext(ProjectContext)
    const {get, loading} = useFetch("/api/")

    useEffect(() => {
        if (projectCtx.project) {
            setSelectedProject(projectCtx.project)
        }

        get("accounts/"+authCtx.user.accountId+"/projects")
            .then(d => setProjects(d))
            .catch(e => console.error(e))
    }, [])

    function handleProjectSelected(e) {
        const projectId = e.target.getAttribute('data-val1')

        if (selectedProject !== null && selectedProject !== undefined) {
            if (selectedProject.id === projectId) {
                setSelectedProject(null)
                return
            }
        }

        const project = projects.find((p) => p.id === projectId)
        setSelectedProject(project)
        projectCtx.setProject(project)
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <Container>
            <Row>
                <Col sm={4}>
                    <ListGroup>
                        {
                            projects.map(project => {
                                return (
                                    <ProjectButton key={project.id} project={project} selectedProject={selectedProject} handleProjectSelected={handleProjectSelected}/>
                                )
                            })
                        }
                    </ListGroup>
                </Col>
                <Col sm={8}>
                    {
                        selectedProject
                            ? <Project project={selectedProject} setSelectedProject={setSelectedProject} setProjects={setProjects}/>
                            : <AccessGuard resourceId={authCtx.user.accountId} action="CreateProject">
                                <CreateProject setProjects={setProjects}/>
                            </AccessGuard>
                    }
                </Col>
            </Row>
        </Container>
    </>
}

function ProjectButton(props) {
    const [readAccess, setReadAccess] = useState(false)
    useAccess(props.project.id, "ReadProject", setReadAccess)

    return <>
        <ListGroupItem
            disabled={!readAccess}
            action
            active={props.selectedProject && props.selectedProject.id === props.project.id}
            onClick={props.handleProjectSelected}
            data-val1={props.project.id}
            className='d-flex justify-content-between align-items-start'
        >
            <div className='ms-2 me-auto'>
                <div className='fw-bold' data-val1={props.project.id}>{props.project.name}</div>
            </div>
        </ListGroupItem>
    </>
}