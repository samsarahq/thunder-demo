import React from 'react';
import { mutate } from 'thunder-react';
import { Trash2 as TrashIcon } from 'react-feather';
import './chat.css'

function deleteMessage(id) {
  mutate({
    query: '{ deleteMessage(id: $id) }',
    variables: { id },
  });
}

function Message({ id, text, username, usernameColor }) {
  return (
    <div className="message">
      <div>
        <div style={{color: usernameColor}} className="message-username">{username}</div>
        <div>{text}</div>
      </div>
      <TrashIcon
        className="message-delete"
        onClick={() => deleteMessage(id)}
        size={16}
      />
    </div>
  )
}
  
class Editor extends React.Component {
  state = { text: '' }

  handleInputChange = (e) => {
    this.setState({text: e.target.value})
  }

  handleSubmit = (e) => {
    mutate({
      query: '{ addMessage(text: $text, sentBy: $sentBy, color: $color) }',
      variables: { text: this.state.text, sentBy: this.props.username, color: this.props.usernameColor },
    }).then(() => {
      this.setState({text: ''});
    });
  }

  handleEnterKey = (e) => {
    if (e.which === 13) {
      this.handleSubmit(e)
    }
  }

  render() {
    return (
      <div className="editor" ref={this.props.setEditorRef}>
        <input
          className="editor-input"
          type="text"
          value={this.state.text}
          onChange={this.handleInputChange}
          onKeyUp={this.handleEnterKey}
        />
        <button
          className="editor-submit button"
          onClick={this.handleSubmit}
        >
          Submit
        </button>
      </div>
    );
  }
}

export default class Chat extends React.Component {
  state = {
    messageContainerHeight: 0
  }

  componentDidMount() {
    this.setState({
      messageContainerHeight: 
        this.chatContainerRef.offsetHeight - 
        this.chatHeaderRef.offsetHeight - 
        this.editorRef.offsetHeight
    }, () => {
      if (this.messagesContainerRef) {
        scrollToBottom(this.messagesContainerRef)
      }
    })
  }

  componentDidUpdate() {
    if (this.messagesContainerRef) {
      scrollToBottom(this.messagesContainerRef)
    }
  }

  render() {
    return (
      <div className="chat-container" ref={el => this.chatContainerRef = el}>
        <div className="chat-header" ref={el => this.chatHeaderRef = el}>
          Let's get chatty
        </div>
        <div
          className="messages-container"
          ref={el => this.messagesContainerRef = el}
          style={{ height: this.state.messageContainerHeight }}
        >
          {this.props.messages.map(props => <Message key={props.id} username={props.sentBy} usernameColor={props.color} {...props} />)}
        </div>
        <Editor setEditorRef={el => this.editorRef = el} username={this.props.username} usernameColor={this.props.usernameColor}/>
      </div>
    )
  }
}

function scrollToBottom(container) {
  if (container.lastElementChild) {
    container.scrollTo({
      behavior: 'smooth',
      top: container.lastElementChild.offsetTop
    });
  }
}

