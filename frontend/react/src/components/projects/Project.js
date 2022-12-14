import useFetch from "../../hooks/useFetch";
import Loader from "../../uikit/Loader";
import Sprints from "./sprints/Sprints";
import {Button, Tab, Tabs} from "react-bootstrap";
import {useContext} from "react";
import {ProjectContext} from "./ProjectContext";

export default function Project(props) {
    const {project, setSelectedProject, setProjects} = props

    if (!project) {
        return <></>
    }

    return <>
        <Tabs defaultActiveKey="sprints">
            <Tab eventKey="sprints" title="Sprints" className="m-2"><Sprints project={project}/></Tab>
            <Tab eventKey="details" title="Details" className="m-2"><ProjectDashboard project={project} setSelectedProject={setSelectedProject} setProjects={setProjects}/></Tab>
        </Tabs>
    </>

}


function ProjectDashboard(props) {
    const {del, loading} = useFetch("/api/")
    const projextCtx = useContext(ProjectContext)

    const {project, setSelectedProject, setProjects} = props

    function handleDeleteClicked() {
        del(`projects/${project.id}`)
            .then(() => {
                setProjects(prev => prev.filter(p => p.id !== project.id))
                setSelectedProject(null)
                projextCtx.setProject(null)
            })
            .catch(e => console.error(e))
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <h1>Project Name</h1>
        <p>{project.name}</p>
        <h1>Description</h1>
        <p>{project.description}</p>

        <Button variant="danger" onClick={handleDeleteClicked}>Delete Project</Button>
    </>
}