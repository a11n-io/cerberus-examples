import {useContext, useState} from "react";
import {AuthContext} from "../../context/AuthContext";
import useFetch from "../../hooks/useFetch";
import Loader from "../../uikit/Loader";
import Input from "../../uikit/Input";
import Btn from "../../uikit/Btn";

export default function CreateProject(props) {
    const auth = useContext(AuthContext)
    const [name, setName] = useState()
    const [description, setDescription] = useState()
    const {post, loading} = useFetch("/api/")

    const {setProjects} = props

    function handleFormSubmit(e) {
        e.preventDefault()
        post("accounts/"+auth.user.accountId+"/projects", {
            name: name,
            description: description
        })
            .then(r => {
                if (r) {
                    setProjects(prev => [...prev, r].sort((a,b) => a.name > b.name))
                }
            })
            .catch(e => console.error(e))
    }

    function handleNameChanged(e) {
        setName(e.target.value)
    }

    function handleDescriptionChanged(e) {
        setDescription(e.target.value)
    }

    if (loading) {
        return <Loader/>
    }

    return <>
        <form onSubmit={handleFormSubmit}>
            <Input required placeholder="Name" onChange={handleNameChanged}/>
            <Input required placeholder="Description" onChange={handleDescriptionChanged}/>
            <Btn type="submit">Create Project</Btn>
        </form>
    </>
}