import {useContext, useEffect, useState} from "react";
import {AuthContext} from "../../context/AuthContext";
import {AccessGuard, Permissions, Roles, Users} from "@a11n-io/cerberus-reactjs";
import "@a11n-io/cerberus-reactjs/dist/index.css"
import {Button, Col, Container, Form, ListGroup, ListGroupItem, Row, Tab, Tabs} from "react-bootstrap";
import Loader from "../../uikit/Loader";
import useFetch from "../../hooks/useFetch";

export default function Settings() {

    return <>
        <Tabs defaultActiveKey='users' className='mb-3'>
            <Tab eventKey='users' title='Users'>
                <Users NoUserSelectedComponent={AddUser}/>
            </Tab>
            <Tab eventKey='roles' title='Roles'>
                <Roles />
            </Tab>
            <Tab eventKey='permissions' title='Permissions'>
                <AccountPermissions />
            </Tab>
        </Tabs>
    </>
}

function AddUser() {
    const authCtx = useContext(AuthContext)
    const { post, loading } = useFetch('/api/')
    const {get} = useFetch(`${process.env.REACT_APP_CERBERUS_API_HOST}/api/`) // get roles from cerberus
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [name, setName] = useState("")
    const [roles, setRoles] = useState([])
    const [selectedRole, setSelectedRole] = useState("")

    useEffect(() => {
        get(`roles?sort=name&order=asc&skip=0&limit=100`, {
            "Authorization": "Bearer " + authCtx.user.cerberusToken
        })
            .then(r => {
                if (r && r.page) {
                    setRoles(r.page)
                } else {
                    setRoles([])
                }
            })
            .catch(e => console.error(e))
    }, [])

    function handleEmailChanged(e) {
        setEmail(e.target.value)
    }

    function handlePasswordChanged(e) {
        setPassword(e.target.value)
    }

    function handleNameChanged(e) {
        setName(e.target.value)
    }

    function handleRoleSelected(e) {
        setSelectedRole(e.target.value)
    }

    function handleFormSubmit(e) {
        e.preventDefault()
        post('users', {
            email: email,
            password: password,
            name: name,
            roleName: selectedRole
        })
            .then((r) => {

            })
            .catch((e) => console.error(e))
    }

    if (loading) {
        return <Loader />
    }

    return <>
        <AccessGuard resourceId={authCtx.user.accountId} action="AddUser">
            <Form onSubmit={handleFormSubmit}>
                <Form.Group className='mb-3'>
                    <Form.Label>Email</Form.Label>
                    <Form.Control
                        type='text'
                        placeholder='Enter email'
                        onChange={handleEmailChanged}
                    />
                </Form.Group>
                <Form.Group className='mb-3'>
                    <Form.Label>Password</Form.Label>
                    <Form.Control
                        type='text'
                        placeholder='Enter password'
                        onChange={handlePasswordChanged}
                    />
                </Form.Group>
                <Form.Group className='mb-3'>
                    <Form.Label>Name</Form.Label>
                    <Form.Control
                        type='text'
                        placeholder='Enter name'
                        onChange={handleNameChanged}
                    />
                </Form.Group>
                <Form.Group className='mb-3'>
                    <Form.Label>Role</Form.Label>
                    <Form.Select onChange={handleRoleSelected}>
                        <option value="">Select a role</option>
                        {
                            roles.map(role => {
                                return (
                                    <option key={role.id} value={role.name}>{role.displayName}</option>
                                )
                            })
                        }
                    </Form.Select>
                </Form.Group>
                <Button disabled={email === "" || password === "" || name === ""
                    || selectedRole === ""
                }
                        variant='primary' type='submit'>
                    Add
                </Button>
            </Form>
        </AccessGuard>
    </>
}

function AccountPermissions() {
    const authCtx = useContext(AuthContext)

    return <>
        <AccessGuard resourceId={authCtx.user.accountId} action="ReadAccountPermissions">
            <Permissions resourceId={authCtx.user.accountId} changeAction="ChangeAccountPermissions"/>
        </AccessGuard>
    </>
}