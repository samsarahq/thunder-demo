import React from 'react';
import { connectGraphQL, mutate } from 'thunder-react';
import { GraphiQLWithFetcher } from './graphiql';

window.puzzleToArray = (puzzle) => {
  let rows = puzzle.split("\n");
  return rows.map((r) => {
    let cols = r.split("|");
    return cols.map((c)=> {
      if (c === ".") {
        return null;
      }
      return parseInt(c, 10);
    });
  });
}

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

let testPuzzle = "6|.|.|.|.|3|.|.|9\n.|.|4|.|1|.|6|.|7\n1|.|.|.|.|.|.|.|.\n.|.|.|4|.|9|.|5|.\n.|2|.|.|.|.|.|7|.\n9|.|7|.|.|.|8|.|4\n.|9|.|.|8|.|.|.|.\n.|8|.|3|.|2|.|9|.\n.|.|.|.|.|.|5|2|.";
console.log(testPuzzle === window.arrayToPuzzle(window.puzzleToArray(testPuzzle)));

let Game = function(props) {
  return (
    <div>team snake</div>
  );
}

// Game = connectGraphQL(Game, () => ({
//   query: ``,
//   variables: {},
//   onlyValidData: true,
// }));

function App() {
  if (window.location.pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else {
    return <Game />;
  }
}

export default App;
