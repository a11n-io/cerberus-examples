import {useContext, useState} from "react";
import useFetch from "../../../../hooks/useFetch";
import Loader from "../../../../uikit/Loader";
import Input from "../../../../uikit/Input";
import Btn from "../../../../uikit/Btn";
import {SprintContext} from "../SprintContext";

export default function CreateStory(props) {
    const sprintCtx = useContext(SprintContext)
    const [description, setDescription] = useState()
    const {post, loading} = useFetch("/api/")

    function handleFormSubmit(e) {
        e.preventDefault()
        post("sprints/"+sprintCtx.sprint.id+"/stories", {
            description: description,
        })
            .then(r => {
                if (r) {
                    props.setStories(prev => [...prev, r].sort((a,b) => a.description > b.description))
                }
            })
            .catch(e => console.error(e))
    }

    function handleDescriptionChanged(e) {
        setDescription(e.target.value)
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <form onSubmit={handleFormSubmit}>
            <Input required placeholder="Description" onChange={handleDescriptionChanged}/>
            <Btn type="submit">Create Story</Btn>
        </form>
    </>
}