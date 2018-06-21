import * as React from 'react';
import { Connection, graphql, GraphQLData, ThunderProvider } from 'thunder-react'; 
import './App.css';

import EventsPage from './components/events_page/EventsPage';
import logo from './logo.svg';

export interface IResult {
  output: string; 
  events: any[]; 
}

const ConnectedEventsPage = graphql<
  GraphQLData<IResult>, 
  IResult
>(
  EventsPage, `query test {
    events(repoID: 1) {
      atMs
      eventId
      repoId
      jsonStr
    }
  }`
);

const connection = new Connection(async () => new WebSocket("ws://localhost:3030/graphql"));

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
        {/* <EventsPage /> */}
        <ThunderProvider connection={connection}>
          <ConnectedEventsPage />
        </ThunderProvider>
      </div>
    );
  }
}

export default App;
