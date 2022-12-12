import {useContext, useEffect, useState} from "react";
import useFetch from "../../../hooks/useFetch";
import Loader from "../../../uikit/Loader";
import {Routes, Route, Link} from "react-router-dom";
import Sprint from "./Sprint";
import CreateSprint from "./CreateSprint";
import {AccessGuard} from "@a11n-io/cerberus-reactjs";
import {ProjectContext} from "../ProjectContext";

export default function Sprints() {
    const projectCtx = useContext(ProjectContext)

    return <>
        <Routes>
            <Route path=":id/*" element={<Sprint/>}/>
            <Route path="*" element={<SprintList project={projectCtx.project}/>}/>
        </Routes>
    </>
}

function SprintList(props) {
    const [sprints, setSprints] = useState([])
    const {get, loading} = useFetch("/api/")
    const [showCreate, setShowCreate] = useState(false)

    const {project} = props

    useEffect(() => {
        get("projects/"+project.id+"/sprints")
            .then(d => {
                if (d) {
                    setSprints(d)
                }
            })
            .catch(e => {
                console.error(e)
                setSprints([])
            })
    }, [project])

    function handleNewClicked(e) {
        e.preventDefault()
        setShowCreate(p => !p)
    }

    if (loading) {
        return <Loader/>
    }

    return <>

        <ul>
            {
                sprints.map(sprint => {
                    return (
                        <li className="nav-item" key={sprint.id}>
                            <AccessGuard
                                resourceId={sprint.id}
                                action="ReadSprint"
                                otherwise={<span>{sprint.sprintNumber}: {sprint.goal}</span>}>
                                <Link to={`/sprints/${sprint.id}`}>
                                    <i>{sprint.sprintNumber}: {sprint.goal}</i>
                                    <i className="m-1">&#8594;</i>
                                </Link>
                            </AccessGuard>
                        </li>
                    )
                })
            }
        </ul>

        <AccessGuard resourceId={project.id} action="CreateSprint">
            {
                !showCreate && <Link to="" onClick={handleNewClicked}>New Sprint</Link>
            }
            {
                showCreate && <CreateSprint
                    sprints={sprints}
                    setSprints={setSprints}
                    setShowCreate={setShowCreate}/>
            }
        </AccessGuard>
    </>
}