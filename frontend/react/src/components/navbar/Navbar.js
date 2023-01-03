import {NavLink, Link, useNavigate} from "react-router-dom";
import {useContext} from "react";
import {AuthContext, AuthGuard} from "../../context/AuthContext";
import {ProjectContext} from "../projects/ProjectContext";
import {SprintContext} from "../projects/sprints/SprintContext";

export default function Navbar() {
    const auth = useContext(AuthContext)
    const projectCtx = useContext(ProjectContext)
    const sprintCtx = useContext(SprintContext)
    const navigate = useNavigate()

    function handleLogout() {
        auth.logout()
        projectCtx.setProject(null)
        sprintCtx.setSprint(null)
        navigate("/login")
    }

    return <>
        <nav className="navbar">
            <NavLink to="/" className="nav-brand">
                Acme Project Manager
            </NavLink>
            <ul>
                <AuthGuard>
                    <NavLinks onLogoutClicked={handleLogout} auth={auth}/>
                </AuthGuard>
            </ul>
        </nav>
    </>
}

function NavLinks(props) {
    const {onLogoutClicked, auth} = props

    return <>
        <li className="nav-item">
            <NavLink to="/projects">Projects</NavLink>
        </li>
        <li className="nav-item">
            <NavLink to="/settings">Settings</NavLink>
        </li>
        <li className="nav-item">
            <Link to="" onClick={onLogoutClicked}>Logout {auth.user.name}</Link>
        </li>
    </>
}