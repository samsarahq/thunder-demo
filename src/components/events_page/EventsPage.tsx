import * as React from "react";

import { get } from "../../api";
import Event from "../event/Event";

import './events_page.css'; 

interface State {
  url?: string;
  events?: any[];
  source?: string; 
}

class EventsPage extends React.Component<{}, State> {
  private timer: NodeJS.Timer; 
  private timerRefresh: number; 

  constructor(props: {}) {
    super(props);
    this.state = {};
    this.timerRefresh = 5000; 
  }

  componentDidMount() {
    this.refreshSource();
    this.timer = setInterval(() => this.refreshSource(), this.timerRefresh);
  }
  
  componentWillUnmount() {
    clearInterval(this.timer); 
  }
  
  refreshSource() {
    // const url = "https://api.github.com/repos/samsarahq/thunder-demo/events";
    // const url = "https://api.github.com/repos/facebookresearch/DensePose/events";
    const url = `https://api.github.com/repos/${this.state.source}/events`;
    get(url).then((json: any[]) => {
      console.log(json);
      this.setState({
        events: json, 
      });
    }); 
  }

  handleInputChange = (event: React.FormEvent<HTMLInputElement>) => {
    this.setState({source: event.currentTarget.value})
  }

  renderEvent(event: any) {
    return <Event event={event} />
  }

  renderPushEvent(event: any) {
    return (
      <div>
        <img src={event.actor.avatar_url} height='25' width='25' />
        <div>{event.actor.display_login}</div>
        <div className='State'>{event.type}</div>
      </div>
    )
  }

  renderEvents() {
    if (!this.state.events) return null;
    return this.state.events.map((event, i) => {
      return (
        <div key={`event-${i}`}>
          {this.renderEvent(event)}
        </div>
      );
    });
  }

  render() {
    return (
      <div className='EventsPage'>
        <input 
          className='EventsPage-source' 
          value={this.state.source} 
          onChange={this.handleInputChange}>
        </input>
        <div className='EventsPage-events'>{this.renderEvents()}</div>
      </div>
    );
  }
}

export default EventsPage;
