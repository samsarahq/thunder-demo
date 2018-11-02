import React from 'react';
import { GraphiQLWithFetcher } from './graphiql';
import { connectGraphQL } from 'thunder-react';
import './app.css'
import SudokuBoard from "./board";
import Chat from './chat'

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

let Sudoku = function(props) {
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

Sudoku = connectGraphQL(Sudoku, (props) => ({
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
  if (pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else if (!isNaN(gameId)) {
    return <Sudoku id={gameId}/>
  } 
  else {
    window.location.pathname = "/1";
    return <Sudoku id={1}/>
  }
}

export default App;
