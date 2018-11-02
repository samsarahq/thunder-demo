import React from 'react';
import { GraphiQLWithFetcher } from './graphiql';
import { connectGraphQL } from 'thunder-react';
import './app.css'
import SudokuBoard from "./board";
import Chat from './chat'

const Sudoku = function(props) {
  const { game, messages } = props.data.value;
  return (
    <div className="app-container">
      <div className="game-container">
        <div className="App-boardWrapper">
          <div className="App-gameName">
            {game.name}
          </div>
          <SudokuBoard id={props.id} initialBoard={game.data} stateBoard={game.state} playerStates={game.playerStates} />
          <div className="u-marginBottomLg" />
        </div>
      </div>
      <Chat messages={messages} />
    </div>
  );
}

const ConnectedSudoku = connectGraphQL(Sudoku, (props) => ({
  query: `
  {
    game(id: $id) {
      data
      state
      solved
      name
      playerStates {
        playerId
        color
        x
        y
      }
    }
    messages {
      id
      text
      sentBy
    }
    currentPlayer {
      color
      name
      playerId
      x
      y
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
  } else if (!isNaN(gameId)) {
    return <ConnectedSudoku id={gameId}/>
  } 
  else {
    window.location.pathname = "/1";
    return <ConnectedSudoku id={1}/>
  }
}

export default App;
