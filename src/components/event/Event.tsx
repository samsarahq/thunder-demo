import * as React from 'react';
import './event.css'; 

export enum EventType {
  WatchEvent,
  ForkEvent,
  PushEvent,
  CreateEvent,
}

interface Props {
  event: any; 
}

interface State {
  
}

type EventInfo = {
  labelText: string; 
}

class Event extends React.Component<Props, State> {

  constructor(props: Props) {
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
        {/* <div className='Event-type'>{event.type}</div> */}
      </div>
    )
  }

  render() {
    let event = this.props.event;

    switch (event.type) {
      case 'PushEvent': {
        return this.renderPushEvent({labelText: 'COMMIT'});
      }
      case 'CreateEvent': {
        return this.renderPushEvent({labelText: 'CREATE'});
      }
      case 'ForkEvent': {
        return this.renderPushEvent({labelText: 'FORK'});
      }
      case 'WatchEvent': {
        return this.renderPushEvent({labelText: 'WATCH'});
      }
      default: {
        return this.renderPushEvent({labelText: 'COMMIT'});
      }
    }
  }
}

export default Event; 
