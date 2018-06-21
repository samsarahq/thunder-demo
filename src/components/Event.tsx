import * as React from 'react';

enum EventType {
  WatchEvent,
  ForkEvent,
  PushEvent,
  CreateEvent,
}

interface Props {
  type: EventType;
  username: string; 
}

interface State {
  
}

class Event extends React.Component<Props, State> {

  constructor(props: Props) {
    super(props);
    this.state = {
      
    }
  }

  public render() {
    return null;
  }
}

export default Event; 
