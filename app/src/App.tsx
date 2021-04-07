import React, { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';
import cookies from 'browser-cookies';

function App() {
  const [greeting, setGreeting] = useState("");

  useEffect(() => {
    async function getGreeting() {
      const token = cookies.get('idToken');
      const response = await fetch("/hello", {
        headers: {
          Authorization: `Bearer ${token}`,
        }
      });
      const body = await response.text();
      setGreeting(body);
    }

    getGreeting();
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <div>
          Greeting: {greeting}
        </div>
      </header>
    </div>
  );
}

export default App;
