import {useContext, useEffect, useState} from "react";
import {AuthContext} from "../../context/AuthContext";
import {Button, Col, Container, Form, ListGroup, ListGroupItem, Row, Tab, Tabs} from "react-bootstrap";
import Loader from "../../uikit/Loader";
import useFetch from "../../hooks/useFetch";

export default function Settings() {

    return <>
        <Tabs defaultActiveKey='users' className='mb-3'>
            <Tab eventKey='users' title='Users'>
                <Users />
            </Tab>
        </Tabs>
    </>
}

function Users() {
    const [users, setUsers] = useState([])
    const [selectedUser, setSelectedUser] = useState(null)
    const {get, loading} = useFetch("/api/")

    useEffect(() => {
        get("users")
            .then(r => setUsers(r))
            .catch(e => console.log(e))
    }, [])


    function handleUserSelected(e) {
        const userId = e.target.getAttribute('data-val1')

        if (selectedUser !== null && selectedUser !== undefined) {
            if (selectedUser.id === userId) {
                setSelectedUser(null)
                return
            }
        }

        const user = users.find((p) => p.id === userId)
        setSelectedUser(user)
    }

    return <>
        <Container>
            <Row>
                <Col sm={4}>
                    <ListGroup>
                        {
                            users.map(user => {
                                return (
                                    <ListGroupItem
                                        key={user.id}
                                        action
                                        active={selectedUser && selectedUser.id === user.id}
                                        onClick={handleUserSelected}
                                        data-val1={user.id}
                                        className='d-flex justify-content-between align-items-start'
                                    >
                                        <div className='ms-2 me-auto'>
                                            <div className='fw-bold' data-val1={user.id}>{user.name}</div>
                                        </div>
                                    </ListGroupItem>
                                )
                            })
                        }
                    </ListGroup>
                </Col>
                <Col sm={8}>
                    {
                        !selectedUser && <AddUser/>
                    }
                </Col>
            </Row>
        </Container>

    </>
}

function AddUser() {
    const authCtx = useContext(AuthContext)
    const { post, loading } = useFetch('/api/')
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [name, setName] = useState("")

    function handleEmailChanged(e) {
        setEmail(e.target.value)
    }

    function handlePasswordChanged(e) {
        setPassword(e.target.value)
    }

    function handleNameChanged(e) {
        setName(e.target.value)
    }

    function handleFormSubmit(e) {
        e.preventDefault()
        post('users', {
            email: email,
            password: password,
            name: name
        })
            .then((r) => {

            })
            .catch((e) => console.error(e))
    }

    if (loading) {
        return <Loader />
    }

    return <>
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

            <Button disabled={email === "" || password === "" || name === ""}
                    variant='primary' type='submit'>
                Add
            </Button>
        </Form>
    </>
}
