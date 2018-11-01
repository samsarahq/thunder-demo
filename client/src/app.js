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
      output += (c == null) ? "." : '' + c;
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
console.log(testPuzzle == window.arrayToPuzzle(window.puzzleToArray(testPuzzle)));


const Editor = React.createClass({
  getInitialState() {
    return {text: ''};
  },

  onSubmit(e) {
    mutate({
      query: '{ addMessage(text: $text) }',
      variables: { text: this.state.text },
    }).then(() => {
      this.setState({text: ''});
    });
  },

  render() {
    return (
      <div>
        <input type="text" value={this.state.text} onChange={e => this.setState({text: e.target.value})} />
        <button onClick={this.onSubmit}>Submit</button>
      </div>
    );
  },
});

function deleteMessage(id) {
  mutate({
    query: '{ deleteMessage(id: $id) }',
    variables: { id },
  });
}

function addReaction(messageId, reaction) {
  mutate({
    query: '{ addReaction(messageId: $messageId, reaction: $reaction) }',
    variables: { messageId, reaction },
  });
}

let Messages = function(props) {
  return (
    <div>
      {props.data.value.messages.map(({id, text, reactions}) =>
        <p key={id}>{text}
          <button onClick={() => deleteMessage(id)}>X</button>
          {reactions.map(({reaction, count}) =>
            <button onClick={() => addReaction(id, reaction)}>{reaction} x{count}</button>
          )}
        </p>
      )}
      <Editor />
    </div>
  );
}
Messages = connectGraphQL(Messages, () => ({
  query: `
  {
    messages {
      id, text
      reactions { reaction count }
    }
  }`,
  variables: {},
  onlyValidData: true,
}));

function App() {
  if (window.location.pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else {
    return <Messages />;
  }
}

export default App;
