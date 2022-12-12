import {useEffect, useState} from "react";
import useFetch from "../../../../hooks/useFetch";
import Loader from "../../../../uikit/Loader";
import {Form, Tab, Tabs} from "react-bootstrap";

export default function Story(props) {
    const {story, setSelectedStory, setStories} = props

    if (!story) {
        return <></>
    }

    return <>
        <Tabs defaultActiveKey="details">
            <Tab eventKey="details" title="Details"><Dashboard story={story} setSelectedStory={setSelectedStory} setStories={setStories}/></Tab>
        </Tabs>
    </>
}


function Dashboard(props) {
    const {get, post, loading} = useFetch("/api/")
    const [users, setUsers] = useState([])
    const [estimate, setEstimate] = useState(0)
    const [status, setStatus] = useState("")
    const [assignee, setAssignee] = useState("")
    const {story, setSelectedStory, setStories} = props

    useEffect(() => {
        get("users")
            .then(d => setUsers(d))
            .catch(e => console.error(e))
    }, [])

    useEffect(() => {
        setEstimate(story.estimation)
        setStatus(story.status)
        setAssignee(story.assignee)
    }, [story])

    function handleEstimateChange(e) {
        setEstimate(e.target.value)
    }

    function handleEstimateBlur(e) {
        post("stories/"+story.id + "/estimate", {
            estimation: estimate
        })
            .then(d => {
                if (d) {
                    setSelectedStory({...d})
                    setStories(prev => [...prev.filter(s => s.id !== story.id), d].sort((a,b) => a.description > b.description))
                }
            })
            .catch(e => console.error(e))
    }

    function handleStatusChange(e) {
        post("stories/"+story.id + "/status", {
            status: e.target.value
        })
            .then(d => {
                if (d) {
                    setSelectedStory({...d})
                    setStories(prev => [...prev.filter(s => s.id !== story.id), d].sort((a,b) => a.description > b.description))
                }
            })
            .catch(e => console.error(e))
    }

    function handleAssigneeChange(e) {
        post("stories/"+story.id + "/assign", {
            userId: e.target.value
        })
            .then(d => {
                if (d) {
                    setSelectedStory({...d})
                    setStories(prev => [...prev.filter(s => s.id !== story.id), d].sort((a,b) => a.description > b.description))
                }
            })
            .catch(e => console.error(e))
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <h2>Description</h2>
        <p>{story.description}</p>
        <Form className="mb-5">
            <Form.Group className="mb-3">
                <Form.Label>Estimate</Form.Label>
                <Form.Control type="number" value={estimate} onChange={handleEstimateChange} onBlur={handleEstimateBlur}/>
            </Form.Group>
            <Form.Group className="mb-3">
                <Form.Label>Status</Form.Label>
                <Form.Select value={status} onChange={handleStatusChange}>
                    <option value="todo">todo</option>
                    <option value="busy">busy</option>
                    <option value="done">done</option>
                </Form.Select>
            </Form.Group>
            <Form.Group className="mb-3">
                <Form.Label>Assignee</Form.Label>
                <Form.Select value={assignee} onChange={handleAssigneeChange}>
                    {
                        users.map(user => {
                            return (
                                <option key={user.id} value={user.id}>{user.name}</option>
                            )
                        })
                    }
                </Form.Select>
            </Form.Group>
        </Form>
    </>
}
