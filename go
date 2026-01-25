<!DOCTYPE html>
<html lang="th">
<head>
<meta charset="UTF-8">
<title>หมากล้อม 19x19 พร้อม Undo/Redo และตัวเลขหมาก</title>
<style>
body { display: flex; flex-direction: column; align-items: center; margin: 0; background: #f0f0f0; font-family: sans-serif; }
canvas { box-shadow: 0 0 10px rgba(0,0,0,0.3); margin-top: 20px; }
#score, #currentPlayer { font-size: 18px; margin-top: 10px; }
button { margin: 5px; padding: 8px 16px; font-size: 16px; cursor: pointer; }
</style>
</head>
<body>

<canvas id="board" width="760" height="760"></canvas>
<div id="score">ดำ: 0 | ขาว: 0</div>
<div id="currentPlayer">ผู้เล่น: ดำ</div>
<div>
  <button id="resetBtn">เริ่มเกมใหม่</button>
  <button id="undoBtn">Undo</button>
  <button id="redoBtn">Redo</button>
  <button id="toggleNumbersBtn">แสดง/ซ่อนตัวเลขหมาก</button>
</div>

<script>
const canvas = document.getElementById('board');
const ctx = canvas.getContext('2d');
const boardSize = 19;
const cellSize = 40;
let currentPlayer = 'black';
let board = Array(boardSize).fill(null).map(() => Array(boardSize).fill(null));
let history = [];
let redoStack = [];
let showNumbers = false;

const scoreDiv = document.getElementById('score');
const playerDiv = document.getElementById('currentPlayer');
const resetBtn = document.getElementById('resetBtn');
const undoBtn = document.getElementById('undoBtn');
const redoBtn = document.getElementById('redoBtn');
const toggleNumbersBtn = document.getElementById('toggleNumbersBtn');

toggleNumbersBtn.addEventListener('click', () => {
  showNumbers = !showNumbers;
  drawBoard();
  drawStones();
});

function drawBoard() {
  ctx.fillStyle = '#f5deb3';
  ctx.fillRect(0, 0, canvas.width, canvas.height);
  ctx.strokeStyle = 'black';
  ctx.lineWidth = 1;
  for (let i = 0; i < boardSize; i++) {
    ctx.beginPath();
    ctx.moveTo(cellSize/2, cellSize/2 + i*cellSize);
    ctx.lineTo(cellSize/2 + cellSize*(boardSize-1), cellSize/2 + i*cellSize);
    ctx.stroke();
    ctx.beginPath();
    ctx.moveTo(cellSize/2 + i*cellSize, cellSize/2);
    ctx.lineTo(cellSize/2 + i*cellSize, cellSize/2 + cellSize*(boardSize-1));
    ctx.stroke();
  }
  const starPoints = [3, 9, 15];
  ctx.fillStyle = 'black';
  starPoints.forEach(r => {
    starPoints.forEach(c => {
      ctx.beginPath();
      ctx.arc(cellSize/2 + c*cellSize, cellSize/2 + r*cellSize, 5, 0, 2*Math.PI);
      ctx.fill();
    });
  });
}

function drawStones() {
  for (let r = 0; r < boardSize; r++) {
    for (let c = 0; c < boardSize; c++) {
      const cell = board[r][c];
      if(cell){
        ctx.beginPath();
        ctx.arc(cellSize/2 + c*cellSize, cellSize/2 + r*cellSize, cellSize*0.4, 0, 2*Math.PI);
        ctx.fillStyle = cell.color;
        ctx.fill();

        if(showNumbers){
          ctx.fillStyle = cell.color === 'black' ? 'white' : 'black';
          ctx.font = `${cellSize*0.4}px sans-serif`;
          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';
          ctx.fillText(cell.moveNumber, cellSize/2 + c*cellSize, cellSize/2 + r*cellSize);
        }
      }
    }
  }
}

function countScore() {
  let blackCount = 0;
  let whiteCount = 0;
  const visited = Array(boardSize).fill(null).map(() => Array(boardSize).fill(false));

  function bfs(r,c){
    const queue = [[r,c]];
    const cells = [];
    let neighbors = new Set();
    while(queue.length){
      const [x,y] = queue.shift();
      if(x<0||x>=boardSize||y<0||y>=boardSize) continue;
      if(visited[x][y]) continue;
      visited[x][y] = true;
      cells.push([x,y]);
      if(board[x][y]===null){
        [[1,0],[0,1],[-1,0],[0,-1]].forEach(([dx,dy])=>queue.push([x+dx,y+dy]));
      }else neighbors.add(board[x][y].color);
    }
    if(neighbors.size===1){
      const color = neighbors.values().next().value;
      cells.forEach(([x,y])=>{
        if(color==='black') blackCount++;
        if(color==='white') whiteCount++;
      });
    }
  }

  for(let r=0;r<boardSize;r++){
    for(let c=0;c<boardSize;c++){
      if(!visited[r][c]&&board[r][c]===null) bfs(r,c);
      else if(board[r][c]&&board[r][c].color==='black') blackCount++;
      else if(board[r][c]&&board[r][c].color==='white') whiteCount++;
    }
  }

  scoreDiv.textContent = `ดำ: ${blackCount} | ขาว: ${whiteCount}`;
}

function updatePlayer(){
  playerDiv.textContent = `ผู้เล่น: ${currentPlayer==='black'?'ดำ':'ขาว'}`;
}

function handleMove(row,col){
  if(row>=0&&row<boardSize&&col>=0&&col<boardSize&&!board[row][col]){
    const moveNumber = history.length + 1;
    history.push(JSON.parse(JSON.stringify(board)));
    redoStack = [];
    board[row][col] = {color: currentPlayer, moveNumber: moveNumber};
    drawBoard();
    drawStones();
    currentPlayer = currentPlayer==='black'?'white':'black';
    countScore();
    updatePlayer();
  }
}

canvas.addEventListener('click', e=>{
  const rect = canvas.getBoundingClientRect();
  const col = Math.round((e.clientX - rect.left - cellSize/2)/cellSize);
  const row = Math.round((e.clientY - rect.top - cellSize/2)/cellSize);
  handleMove(row,col);
});

resetBtn.addEventListener('click', ()=>{
  history.push(JSON.parse(JSON.stringify(board)));
  redoStack=[];
  board=Array(boardSize).fill(null).map(()=>Array(boardSize).fill(null));
  currentPlayer='black';
  drawBoard();
  drawStones();
  countScore();
  updatePlayer();
});

undoBtn.addEventListener('click', ()=>{
  if(history.length){
    redoStack.push(JSON.parse(JSON.stringify(board)));
    board = history.pop();
    drawBoard();
    drawStones();
    countScore();
    updatePlayer();
  }
});

redoBtn.addEventListener('click', ()=>{
  if(redoStack.length){
    history.push(JSON.parse(JSON.stringify(board)));
    board = redoStack.pop();
    drawBoard();
    drawStones();
    countScore();
    updatePlayer();
  }
});

// เริ่มต้น
drawBoard();
drawStones();
countScore();
updatePlayer();
</script>

</body>
</html>
