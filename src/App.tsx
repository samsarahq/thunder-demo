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
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          To get started, edit <code>src/App.tsx</code> and save to reload.
        </p>
        <EventsPage />
      </div>
    );
  }
}

export default App;
