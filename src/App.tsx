import * as React from 'react';
import './App.css';

import EventsPage from './components/events_page/EventsPage';
import logo from './logo.svg';

class App extends React.Component {

  public render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to Thunder</h1>
        </header>
        <p className="App-intro">
          Enter any repo in the form <code>owner/repo_name</code> to watch its activities!
        </p>
        <EventsPage />
      </div>
    );
  }
}

export default App;
