// sudoku validator adopted from https://gist.github.com/0xsven/5cea8fc8c7e794554419

export const getValidCells = (board) => {

    let _rows, _cols, _grid, validCells;

    board = board.map((row, r) => {
        return row.map((col, c) => {
            return {
                val: col,
                row: r,
                col: c
            }
        });
    });


    // reorganize board into three structures
    _rows = board;
    _cols = [];
    _grid = [];
    validCells = [];

    // Prefilling the structures with empty array objects
    for (var i = 0; i < 9; i++) {
        _cols.push([]);
        _grid.push([]);
        validCells.push([]);
    }
    
    for (var row = 0; row < 9; row++) {
        for (var col = 0; col < 9; col++) {
            validCells[row][col] = true;
            // Save each column in a new row
            _cols[col][row] = board[row][col];

            // Calculate grid identifiers
            const gridRow = Math.floor( row / 3 );
            const gridCol = Math.floor( col / 3 );
            const gridIndex = gridRow * 3 + gridCol;

            // Save each grid in a new row
            _grid[gridIndex].push(board[row][col]);       
        }
    }

    // validate rows 
    const _validate = function(data){
        for (let row = 0; row < 9; row++) {
            data[row].sort((a, b) => a.val && b.val ? a.val - b.val : 0);
            for (let col = 0; col < 9; col++) {
                let cell = data[row][col], next_cell = data[row][col + 1];
                
                if (cell.val === null || cell.val === "") {
                    continue;
                }

                // check if numbers are unique
                if (col !== 8 && cell.val === next_cell.val){
                    validCells[cell.row][cell.col] = false;
                    validCells[next_cell.row][next_cell.col] = false;
                }
            }
        }
    };

    _validate(_rows);
    _validate(_cols);
    _validate(_grid);

    return validCells;
}