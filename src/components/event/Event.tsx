import * as React from 'react';
import './event.css'; 

export enum EventType {
  WatchEvent,
  ForkEvent,
  PushEvent,
  CreateEvent,
}

export interface EventProps {
  event: {
    actor: {
      avatar_url: string; 
      display_login: string; 
    }
    commit_msg: string; 
    url: string; 
  } 
  index: number; 
}

interface State {
  
}

type EventInfo = {
  labelText: string; 
}

class Event extends React.Component<EventProps, State> {

  constructor(props: EventProps) {
    super(props);
    this.state = {
      
    }
  }

  renderPushEvent(info: EventInfo) {
    let event = this.props.event; 
    return (
      <a href={event.url} target='_blank'>
        <div className='Event' style={{animationDelay: `${this.props.index*100}ms`}}>
        <div className='Event-labelText'>{info.labelText}</div>
        <div className='Event-userContainer'>
          <img className='Event-userAvatar' src={event.actor.avatar_url} height='25' width='25' />
          <div className='Event-userLogin'>{event.actor.display_login}</div>
        </div>
        <div className='Event-commitMsg'>{event.commit_msg}</div>
      </div>
      </a>
    )
  }

  render() {
    return this.renderPushEvent({labelText: 'COMMIT'});
  }
}

export default Event; 
