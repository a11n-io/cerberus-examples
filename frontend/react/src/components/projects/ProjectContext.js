import {createContext, useState} from "react";
import useSessionStorageState from 'use-session-storage-state';

const ProjectContext = createContext(null)

function ProjectProvider(props) {
    const [project, setProject] = useSessionStorageState(`acme-project`, {defaultValue: null});

    const value = {
        project: project,
        setProject: setProject
    }

    return (
        <ProjectContext.Provider value={value}>
            {props.children}
        </ProjectContext.Provider>
    )
}

export {ProjectContext, ProjectProvider}