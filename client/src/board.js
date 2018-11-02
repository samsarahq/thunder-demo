import React from 'react';
import classNames from 'classnames';
import { connectGraphQL, mutate } from 'thunder-react';

const WIDTH = 9;
const HEIGHT = 9;
const board =  Array.from(Array(WIDTH), _ => Array(HEIGHT).fill(0))
 
export default class SudokuBoard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      x: 0,
      y: 0,
    }
  }

  puzzleToArray = (puzzle) => {
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
  
  componentWillMount() {
    window.onkeydown = this.handleKeyDown
  }

  handleCellChange = (x,y,val) => {
    const cellValue = parseInt(val, 10)
    if(cellValue && 1 <= cellValue && cellValue <= 9) {
      mutate({
        query: `{updateGame(id:$id, col: $col, row: $row, val: $val)}`,
        variables: {
          id: this.props.id,
          col: x,
        row: y,
      val: cellValue}
      })
    }
  }

  handleKeyDown = (event) => {
    let {x, y} = this.state
    switch(event.key) {
      case "ArrowUp": y = Math.max(y - 1, 0); break;
      case "ArrowDown": y = Math.min(y + 1, HEIGHT - 1); break;
      case "ArrowLeft": x = Math.max(x - 1, 0); break;
      case "ArrowRight": x = Math.min(x + 1, WIDTH - 1); break;
      default:
        this.handleCellChange(x, y,event.key);
    }
    console.log(board, board[0][1])
    this.setState({x, y})
  }

  handleClick = (x, y) => () => {
    this.setState({x, y})
  }

  render() {
    const board = this.puzzleToArray(this.props.board)
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
