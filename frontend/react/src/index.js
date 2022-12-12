import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import {AuthProvider} from "./context/AuthContext";
import {CerberusProvider} from "@a11n-io/cerberus-reactjs";

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
      <CerberusProvider apiHost={process.env.REACT_APP_CERBERUS_API_HOST} socketHost={process.env.REACT_APP_CERBERUS_WS_HOST}>
          <AuthProvider>
              <App />
          </AuthProvider>
      </CerberusProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
