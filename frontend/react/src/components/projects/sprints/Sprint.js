import {SprintContext} from "./SprintContext";
import {Route, Routes, useParams} from "react-router-dom";
import {useContext, useEffect} from "react";
import useFetch from "../../../hooks/useFetch";
import Loader from "../../../uikit/Loader";
import Btn from "../../../uikit/Btn";
import Stories from "./stories/Stories";
import {AccessGuard, Permissions} from "@a11n-io/cerberus-reactjs";
import {Card, Tab, Tabs} from "react-bootstrap";
import {ProjectContext} from "../ProjectContext";

export default function Sprint() {
    const params = useParams()
    const projectCtx = useContext(ProjectContext)
    const sprintCtx = useContext(SprintContext)
    const {get, loading} = useFetch("/api/")

    useEffect(() => {
        get("sprints/"+params.id)
            .then(d => sprintCtx.setSprint(d))
            .catch(e => console.error(e))
    }, [])

    if (loading) {
        return <Loader/>
    }

    if (!sprintCtx.sprint) {
        return <>Could not get sprint</>
    }

    return <>
        <h1>Project: {projectCtx.project.name}</h1>
        <Routes>
            <Route path="*" element={<Dashboard/>}/>
        </Routes>
    </>
}


function Dashboard() {
    const sprintCtx = useContext(SprintContext)
    const sprint = sprintCtx.sprint

    return <>
        <Card className="m-2">
            <Card.Header>Sprint {sprint.sprintNumber}: {sprint.goal}</Card.Header>
            <Card.Body>
                <Tabs defaultActiveKey="stories">
                    <Tab eventKey="stories" title="Stories" className="m-2"><Stories /></Tab>
                    <Tab eventKey="details" title="Details" className="m-2"><Details /></Tab>
                    <Tab eventKey="permissions" title="Permissions" className="m-2"><SprintPermissions /></Tab>
                </Tabs>
            </Card.Body>
        </Card>

    </>
}

function Details() {
    const sprintCtx = useContext(SprintContext)
    const sprint = sprintCtx.sprint

    return <>
        <h2>Goal</h2>
        <p>{sprint.goal}</p>
        {
            sprint.startDate === 0 ?
                <ChangeSprint sprintCtx={sprintCtx} start={true}/> :
                <>
                    <p>Started on {new Date(sprint.startDate * 1000).toDateString()}</p>
                    {
                        sprint.endDate === 0 ?
                            <ChangeSprint sprintCtx={sprintCtx} start={false}/> :
                            <>
                                <p>Ended on {new Date(sprint.endDate * 1000).toDateString()}</p>
                            </>
                    }
                </>
        }
    </>
}

function ChangeSprint(props) {
    const {post, loading} = useFetch("/api/")
    const {sprintCtx, start} = props

    function handleButtonClicked() {
        post("sprints/"+sprintCtx.sprint.id+"/" + (start ? "start" : "end"))
            .then(d => {
                if (d) {
                    sprintCtx.setSprint(d)
                }
            })
            .catch(e => console.error(e))
    }

    return <>
        <AccessGuard resourceId={sprintCtx.sprint.id} action="StartSprint">
            <Btn onClick={handleButtonClicked}>{start ? "Start" : "End"} sprint</Btn>
        </AccessGuard>
    </>
}

function SprintPermissions() {
    const sprintCtx = useContext(SprintContext)

    return <>
        <AccessGuard resourceId={sprintCtx.sprint.id} action="ReadSprintPermissions">
            <Permissions resourceId={sprintCtx.sprint.id} changeAction="ChangeSprintPermissions"/>
        </AccessGuard>
    </>
}