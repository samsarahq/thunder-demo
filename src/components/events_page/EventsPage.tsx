import * as React from "react";

import { get } from "../../api";
import Event from "../event/Event";

import './events_page.css'; 

interface State {
  url?: string;
  events?: any[];
}

class EventsPage extends React.Component<{}, State> {
  private timer: NodeJS.Timer; 
  private timerRefresh: number; 

  constructor(props: {}) {
    super(props);
    this.state = {};
    this.timerRefresh = 10000; 
  }

  public componentDidMount() {
    this.refreshSource();
    this.timer = setInterval(() => this.refreshSource(), this.timerRefresh);
  }
  
  public componentWillUnmount() {
    clearInterval(this.timer); 
  }
  
  public refreshSource() {
    const url = "https://api.github.com/repos/samsarahq/thunder-demo/events";
    get(url).then((json: any[]) => {
      console.log(json);
      this.setState({
        events: json
      });
    }); 
  }

  public renderEvent(event: any) {
    return <Event event={event} />
  }

  public renderPushEvent(event: any) {
    return (
      <div>
        <img src={event.actor.avatar_url} height='25' width='25' />
        <div>{event.actor.display_login}</div>
        <div className='State'>{event.type}</div>
      </div>
    )
  }

  public renderEvents() {
    if (!this.state.events) return null;
    return this.state.events.map((event, i) => {
      return (
        <div key={`event-${i}`}>
          {this.renderEvent(event)}
        </div>
      );
    });
  }

  public render() {
    return <div className='EventsPage'>{this.renderEvents()}</div>;
  }
}

export default EventsPage;
