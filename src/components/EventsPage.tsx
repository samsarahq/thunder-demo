import * as React from "react";

import { get } from "../api";

interface State {
  url?: string;
  events?: any[];
}

class EventsPage extends React.Component<{}, State> {
  constructor(props: {}) {
    super(props);
    this.state = {};
  }

  public componentDidMount() {
    const url = "https://api.github.com/repos/lbkchen/thunder-demo/events";
    get(url).then((json: any[]) => {
      console.log(json);
      this.setState({
        events: json
      });
    });
  }

  public renderEvent(event: any) {
    switch (event.type) {
      case 'PushEvent': {
        return this.renderPushEvent(event);
      }
      default: {
        return this.renderPushEvent(event);
      }
    }
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
    return <div className='container'>{this.renderEvents()}</div>;
  }
}

export default EventsPage;
