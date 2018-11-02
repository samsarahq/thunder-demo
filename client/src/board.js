import React from 'react';
import classNames from 'classnames';
import { connectGraphQL } from 'thunder-react';

const WIDTH = 9;
const HEIGHT = 9;
const board = Array.from(Array(WIDTH), _ => Array(HEIGHT).fill(0));

class SudokuBoard extends React.Component {
  constructor() {
    super()
    this.state = {
      x: 0,
      y: 0,
    }
  }
  
  componentWillMount() {
    window.onkeydown = this.handleKeyDown
  }

  handleKeyDown = (event) => {
    let {x, y} = this.state
    switch(event.key) {
      case "ArrowUp": y = Math.max(y - 1, 0); break;
      case "ArrowDown": y = Math.min(y + 1, HEIGHT - 1); break;
      case "ArrowLeft": x = Math.max(x - 1, 0); break;
      case "ArrowRight": x = Math.min(x + 1, WIDTH - 1); break;
      default:
        const cellValue = parseInt(event.key, 10)
        if(cellValue && 1 <= cellValue && cellValue <= 9) {
          board[y][x] = cellValue
        }
    }
    console.log(board, board[0][1])
    this.setState({x, y})
  }

  handleClick = (x, y) => () => {
    this.setState({x, y})
  }

  render() {
    return <div>{board.map(
      (row, i) => <div key={i}>{
        row.map((cell, j) => <BoardCell value={cell} isSelected={Boolean(i === this.state.y && j === this.state.x)} onClick={this.handleClick(j, i)} key={j} />)
      }</div>)
    }</div>
  }
}

const focusRef = isSelected => ref => isSelected && ref && console.log("isSelected", isSelected, ref, ref.focus()) && ref.focus();

const BoardCell = (props) => {return (
  <div className={classNames("BoardCell", {"is-selected": props.isSelected})} onClick={props.onClick} ref={focusRef(props.isSelected)}>{props.value ? props.value : "" }</div>
)}

export default connectGraphQL(SudokuBoard, () => ({
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
