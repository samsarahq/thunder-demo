import React from 'react';
import { GraphiQLWithFetcher } from './graphiql';
import { connectGraphQL, mutate } from 'thunder-react';
import './app.css'
import SudokuBoard from "./board";

window.arrayToPuzzle = (arr) => {
  let output = "";
  arr.forEach((r, i) => {
    r.forEach((c, j) => {
      output += (c === null) ? "." : '' + c;
      if (j < r.length - 1)
        output += '|';
    });
    if (i < arr.length - 1) {
      output += "\n";
    }
  });
  return output
}

//let testPuzzle = "6|.|.|.|.|3|.|.|9\n.|.|4|.|1|.|6|.|7\n1|.|.|.|.|.|.|.|.\n.|.|.|4|.|9|.|5|.\n.|2|.|.|.|.|.|7|.\n9|.|7|.|.|.|8|.|4\n.|9|.|.|8|.|.|.|.\n.|8|.|3|.|2|.|9|.\n.|.|.|.|.|.|5|2|.";
//console.log(testPuzzle === window.arrayToPuzzle(window.puzzleToArray(testPuzzle)));

class Editor extends React.Component {
  state = { text: '' }

  handleInputChange = (e) => {
    this.setState({text: e.target.value})
  }

  handleSubmit = (e) => {
    mutate({
      query: '{ addMessage(text: $text) }',
      variables: { text: this.state.text },
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
      <div className="editor">
        <input
          className="editor-input"
          type="text"
          value={this.state.text}
          onChange={this.handleInputChange}
          onKeyUp={this.handleEnterKey}
        />
        <button
          className="editor-submit"
          onClick={this.handleSubmit}
        >
          Submit
        </button>
      </div>
    );
  }
}

function deleteMessage(id) {
  mutate({
    query: '{ deleteMessage(id: $id) }',
    variables: { id },
  });
}

function Message({ id, text, username }) {
  return (
    <div className="message">
      <div>
        <div>{username}</div>
        <div>{text}</div>
      </div>
      <button className="message-delete" onClick={() => deleteMessage(id)}>
        X
      </button>
    </div>
  )
}

function Chat({ messages }) {
  return (
    <div>
      {messages.map(props => <Message key={props.id} username="user" {...props} />)}
    </div>
  )
}

let Sudoku = function(props) {
  const { game, messages } = props.data.value;
  return (
    <div className="app-container">
      <div className="game-container">
      <div className="App-boardWrapper">
        <div className="App-gameName">
          {game.name}
        </div>
        <SudokuBoard id={props.id} board={game.data} playerStates={game.playerStates} />
        <div className="u-marginBottomLg" />
      </div>
      </div>
      <div className="chat-container">
        <Chat messages={messages} />
        <Editor />
      </div>
    </div>
  );
}

Sudoku = connectGraphQL(Sudoku, (props) => ({
  query: `
  {
    game(id: $id) {
      data
      name
      playerStates {
        playerId
        color
        x
        y
      }
    }
    messages {
      id, text
    }
  }`,
  variables: {
    id: props.id
  },
  onlyValidData: true,
}));

function App() {
  const pathname = window.location.pathname;
  const gameId = parseInt(pathname.slice(1),10);
  console.log(gameId);
  if (pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else if (!isNaN(gameId)) {
    return <Sudoku id={gameId}/>
  } 
  else {
    return <Sudoku id={1}/>
  }
}

export default App;
