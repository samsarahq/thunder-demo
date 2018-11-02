import React from 'react';
import classNames from 'classnames';
import { mutate } from 'thunder-react';
import {getValidCells} from './sudoku';
import './board.css'

const WIDTH = 9;
const HEIGHT = 9;
const CurrentPlayerState = {
  color: "#448AFF",
};

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
    document.addEventListener('keydown',this.handleKeyDown)
  }

  handleCellChange = (x,y,val) => {
    let cellValue = parseInt(val, 10);
    if (val === "Backspace") {
      cellValue = 0;
    }
    if(!isNaN(cellValue) && 0 <= cellValue && cellValue <= 9) {
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
    console.log({ event })
    let {x, y} = this.state
    switch(event.key) {
      case "ArrowUp": y = Math.max(y - 1, 0); break;
      case "ArrowDown": y = Math.min(y + 1, HEIGHT - 1); break;
      case "ArrowLeft": x = Math.max(x - 1, 0); break;
      case "ArrowRight": x = Math.min(x + 1, WIDTH - 1); break;
      default:
        this.handleCellChange(x, y,event.key);
    }
    this.setState({x, y})
  }

  handleClick = (x, y) => () => {
    this.setState({x, y})
  }

  getPlayerState = (x, y) => {
    if(x === this.state.x && y === this.state.y) {
      return CurrentPlayerState
    }
    return this.props.playerStates.find(p => p.x === x && p.y === y)
  }

  render() {
    const stateBoard = this.puzzleToArray(this.props.stateBoard)
    const initialBoard = this.puzzleToArray(this.props.initialBoard)
    const validBoard = getValidCells(stateBoard);
    return <div className="board">{stateBoard.map(
      (row, i) => <div className="row" key={i}>{
        row.map((cell, j) => 
          <BoardCell disabled={initialBoard[i][j]!== null} valid={validBoard[i][j]}value={cell} playerState={this.getPlayerState(j, i)} onClick={this.handleClick(j, i)} key={j} />
        )
      }</div>)
    }</div>
  }
}

const focusRef = isSelected => ref => isSelected && ref && ref.focus();

const BoardCell = (props) => {
  const style = props.playerState && {
    outline: `3px solid ${props.playerState.color}`,
    boxShadow: `0 0 0 3px ${props.playerState.color}`
  };

  const disabledStyle = props.disabled && {
    background: '#bdc3c7'
  };

  const invalidStyle = !props.valid && {
    background: '#e74c3c'
  }

  return (
    <div className={classNames("BoardCell", {"is-selected": Boolean(props.playerState)})} style={disabledStyle || invalidStyle || style} onClick={!props.disabled && props.onClick} ref={focusRef(props.isSelected)}>{props.value ? props.value : "" }</div>
  )
}
