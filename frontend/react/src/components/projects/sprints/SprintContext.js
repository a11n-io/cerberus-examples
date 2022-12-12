import {createContext, useState} from "react";
import useSessionStorageState from 'use-session-storage-state';

const SprintContext = createContext(null)

function SprintProvider(props) {
    const [sprint, setSprint] = useSessionStorageState(`acme-sprint`, {defaultValue: null});

    const value = {
        sprint: sprint,
        setSprint: setSprint
    }

    return (
        <SprintContext.Provider value={value}>
            {props.children}
        </SprintContext.Provider>
    )
}

export {SprintContext, SprintProvider}