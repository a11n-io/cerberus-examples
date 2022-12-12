import {useContext, useState} from "react";
import {ProjectContext} from "../ProjectContext";
import useFetch from "../../../hooks/useFetch";
import Loader from "../../../uikit/Loader";
import Input from "../../../uikit/Input";
import Btn from "../../../uikit/Btn";

export default function CreateSprint(props) {
    const projectCtx = useContext(ProjectContext)
    const [goal, setGoal] = useState()
    const {post, loading} = useFetch("/api/")

    function handleFormSubmit(e) {
        e.preventDefault()
        post("projects/"+projectCtx.project.id+"/sprints", {
            goal: goal,
        })
            .then(r => {
                if (r) {
                    props.setSprints([...props.sprints, r])
                    props.setShowCreate(false)
                }
            })
            .catch(e => console.error(e))
    }

    function handleGoalChanged(e) {
        setGoal(e.target.value)
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <form onSubmit={handleFormSubmit}>
            <Input required placeholder="Goal" onChange={handleGoalChanged}/>
            <Btn type="submit">Create Sprint</Btn>
        </form>
    </>
}