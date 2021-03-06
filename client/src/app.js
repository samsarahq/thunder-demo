import React from 'react';
import { GraphiQLWithFetcher } from './graphiql';
import { connectGraphQL } from 'thunder-react';
import './app.css'
import SudokuBoard from "./board";
import NewGame from "./new_game";
import Chat from './chat'

const Sudoku = function(props) {
  const { game, messages, currentPlayer } = props.data.value;
  return (
    <div className="app-container">
      <div className="game-container">
        <div className="App-boardWrapper">
          <div className="App-gameName">
            {game.name}
          </div>
          <SudokuBoard id={props.id} currentPlayer={currentPlayer} initialBoard={game.data} stateBoard={game.state} playerStates={game.players} />
          <div className="u-marginBottomLg" />
        </div>
      </div>
      <Chat messages={messages} username={currentPlayer.name} usernameColor={currentPlayer.color}/>
    </div>
  );
}

const ConnectedSudoku = connectGraphQL(Sudoku, (props) => ({
  query: `
  {
    currentPlayer {
      id
      name
      color
      x
      y
    }
    game(id: $id) {
      data
      state
      solved
      name
      players {
        id
        name
        color
        x
        y
      }
    }
    messages {
      id
      text
      sentBy
      color
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
  if (pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else if (pathname === "/new") {
    return <NewGame/>
  } else if (!isNaN(gameId)) {
    return <ConnectedSudoku id={gameId}/>
  } 
  else {
    window.location.pathname = "/1";
    return <ConnectedSudoku id={1}/>
  }
}

export default App;
