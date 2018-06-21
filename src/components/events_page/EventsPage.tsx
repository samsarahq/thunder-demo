import * as React from "react";
import { GraphQLData } from 'thunder-react'; 

import { get } from "../../api";
import Event from "../event/Event";
import { IResult } from "../../App";

import './events_page.css'; 

interface State {
  url?: string;
  events?: any[];
  source: string; 
}

class EventsPage extends React.Component<GraphQLData<IResult>, State> {
  private timer: NodeJS.Timer; 
  private timerRefresh: number; 
  private prevSource: string;
  private stallCount: number; 

  constructor(props: GraphQLData<IResult>) {
    super(props);
    this.state = {
      source: 'samsarahq/thunder-demo', 
    };
    this.timerRefresh = 100000; 
    this.prevSource = ''; 
    this.stallCount = 0; 
  }

  componentDidMount() {
    this.refreshSource();
    this.timer = setInterval(() => this.refreshSource(), this.timerRefresh);
  }
  
  componentWillUnmount() {
    clearInterval(this.timer); 
  }
  
  refreshSource() {
    if (this.state.source === this.prevSource) {
      this.stallCount += 1; 
      if (this.stallCount < 5) {
        return;
      } else {
        this.stallCount = 0; 
      }
    }

    const url = `https://api.github.com/repos/${this.state.source}/events`;
    get(url).then((json: any[]) => {
      console.log(json);
      this.setState({
        events: json, 
      });
    }); 
  }

  handleInputChange = (event: React.FormEvent<HTMLInputElement>) => {
    this.prevSource = this.state.source; 
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
