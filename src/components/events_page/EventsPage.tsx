import * as React from "react";
import { GraphQLData } from 'thunder-react'; 

import Event from "../event/Event";
import { IResult } from "../../App";

import './events_page.css'; 

interface State {
  url?: string;
  events?: any[];
  source: string; 
}

class EventsPage extends React.Component<GraphQLData<IResult>, State> {

  constructor(props: GraphQLData<IResult>) {
    super(props); 
    this.state = {
      source: 'samsarahq/thunder-demo', 
    };
  }
  
  handleInputChange = (event: React.FormEvent<HTMLInputElement>) => {
    this.setState({source: event.currentTarget.value})
  }

  getEvent(event: any) {
    const eventObj = JSON.parse(event.jsonStr); 
    console.log(eventObj); 
    return {
      actor: {
        avatar_url: eventObj.author.avatar_url, 
        display_login: eventObj.author.login, 
      }
    };
  }

  renderEvent(event: any) {
    return <Event event={this.getEvent(event)} />
  }

  renderEvents(events: any[]) {
    if (!events) return null;
    return events.map((event, i) => {
      return (
        <div key={`event-${i}`}>
          {this.renderEvent(event)}
        </div>
      );
    });
  }

  render() {
    if (!this.props.data.value) {
      return null; 
    }
    const events = this.props.data.value.events; 
    return (
      <div className='EventsPage'>
        <input 
          className='EventsPage-source' 
          value={this.state.source} 
          onChange={this.handleInputChange}>
        </input>
        <div className='EventsPage-events'>{this.renderEvents(events)}</div>
      </div>
    );
  }
}

export default EventsPage;
