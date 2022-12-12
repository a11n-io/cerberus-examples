import {useContext, useEffect, useState} from "react";
import useFetch from "../../../../hooks/useFetch";
import Loader from "../../../../uikit/Loader";
import {Routes, Route} from "react-router-dom";
import CreateStory from "./CreateStory";
import Story from "./Story";
import {SprintContext} from "../SprintContext";
import {AccessGuard, useAccess} from "@a11n-io/cerberus-reactjs";
import {Col, Container, ListGroup, ListGroupItem, Row} from "react-bootstrap";

export default function Stories() {
    return <>
        <Routes>
            <Route exact path="/" element={<StoryList/>}/>
        </Routes>
    </>
}

function StoryList() {
    const [stories, setStories] = useState([])
    const [selectedStory, setSelectedStory] = useState(null)
    const sprintCtx = useContext(SprintContext)
    const {get, loading} = useFetch("/api/")

    useEffect(() => {
        get("sprints/"+sprintCtx.sprint.id+"/stories")
            .then(d => {
                if (d) {
                    setStories(d)
                }
            })
            .catch(e => console.error(e))
    }, [])

    function handleStorySelected(e) {
        const storyId = e.target.getAttribute('data-val1')

        if (selectedStory !== null && selectedStory !== undefined) {
            if (selectedStory.id === storyId) {
                setSelectedStory(null)
                return
            }
        }

        setSelectedStory(stories.find((s) => s.id === storyId))
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
                            stories.map(story => {
                                return (
                                    <StoryButton key={story.id} story={story} selectedStory={selectedStory} handleStorySelected={handleStorySelected}/>
                                )
                            })
                        }
                    </ListGroup>
                </Col>
                <Col sm={8}>
                    {
                        selectedStory
                            ? <Story story={selectedStory} setSelectedStory={setSelectedStory} setStories={setStories}/>
                            : <AccessGuard resourceId={sprintCtx.sprint.id} action="CreateStory">
                                    <CreateStory setStories={setStories}/>
                              </AccessGuard>
                    }
                </Col>
            </Row>
        </Container>
    </>
}

function StoryButton(props) {
    const [readAccess, setReadAccess] = useState(false)
    useAccess(props.story.id, "ReadStory", setReadAccess)

    return <>
        <ListGroupItem
            disabled={!readAccess}
            action
            active={props.selectedStory && props.selectedStory.id === props.story.id}
            onClick={props.handleStorySelected}
            data-val1={props.story.id}
            className='d-flex justify-content-between align-items-start'
        >
            <div className='ms-2 me-auto'>
                <div className='fw-bold' data-val1={props.story.id}>{props.story.description}</div>
            </div>
        </ListGroupItem>
    </>
}