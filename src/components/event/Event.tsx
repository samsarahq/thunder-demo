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
  } 
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
      <div className='Event'>
        <div className='Event-labelText'>{info.labelText}</div>
        <div className='Event-userContainer'>
          <img className='Event-userAvatar' src={event.actor.avatar_url} height='25' width='25' />
          <div className='Event-userLogin'>{event.actor.display_login}</div>
        </div>
      </div>
    )
  }

  render() {
    return this.renderPushEvent({labelText: 'COMMIT'});
  }
}

export default Event; 
