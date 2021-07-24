import './App.css';
import React, { useState } from "react";
// import getGame from'./gameManager.js';

const OneSecondInMilliSeconds = 1000

async function getGame(friendlyID) {
  const requestPath = `${window.location.href}GameState?friendlyID=${friendlyID}`
  console.log("Getting state:: " + requestPath)

  const myRequest = new Request(requestPath);

  var game = await fetch(myRequest)
    .then(response => {
      if (response.status === 200) {
        return response.json()
      } else {
        throw new Error('Request failed');
      }
    })
    .then(data => {
      console.log(data)
      return data;
    })
    .catch(console.error);

    return game
}

function GameStatus(props) {

  let currentplayers = [];
  if (props.game !== undefined && props.game.State !== undefined && props.game.State.Players.length > 0) {
    for (let player of props.game.State.Players) {
      let DiceString = ""
      if (player.Dice !== null) {
        DiceString = player.Dice.toString();
      }
      currentplayers.push(<p style={{textAlign:"left"}}>Name: {player.Name} Remaining Dice: {player.RemainingDiceCount} Dice: {DiceString}</p>)
    }
  }

  if (currentplayers.length > 0) {
    return (
      <div id="gameState" style={{ fontSize: 14, textAlign: "center", backgroundColor: "#40F2FF" }}>
        <p>Current Bet Value: {props.game.State.CurrentBet.Value}</p>
        <p>Current Bet Count: {props.game.State.CurrentBet.Count}</p>
        <p>Current Player (i - improve this to return playername (as a new attr?)): {props.game.State.CurrentPlayer}</p>
        {currentplayers}
        <p>Active round: {props.game.State.RoundActive.toString()}</p>
        <p>Round: {props.game.State.Round}</p>
        <p>Started: {props.game.State.Started.toString()}</p>

      </div>
    );
  } else {
    return null
  }
}

const emptyGame = {
  ID: "",
  FriendlyID: "",
  State: {
    CurrentBet: { Count: 0, Value: 0 },
    CurrentPlayer: 0,
    Players: [],
    Round: 0,
    RoundActive: false,
    Started: false,
  },
}

function App() {

  const [game, setGame] = useState(emptyGame);

  React.useEffect(() => {
    getGame("Blam").then(data => setGame(data));
  }, []);

  return (
    <div className="App">
      {/* Make this a component */}
      <form id="newGame" style={{ fontSize: 14, textAlign: "center", backgroundColor: "#71B2EB" }}>
        <label htmlFor="CreatorName">Creator Name:</label>
        <input type="text" id="CreatorName" name="CreatorName" defaultValue="Blam" />

        <label htmlFor="GameCode">Game Code:</label>
        <input type="text" id="GameCode" name="GameCode" defaultValue="Blam" />

        <label htmlFor="Playercount">Player count:</label>
        <input type="text" id="Playerccount" name="Playercount" defaultValue="Blam" />

        <label htmlFor="Diceperplayer">Dice per player:</label>
        <input type="text" id="Diceperplayer" name="Diceperplayer" defaultValue="Blam" />

        <input type="submit" value="Create New Game" />
      </form>

      <form id="joinGame" style={{ fontSize: 14, textAlign: "center", backgroundColor: "#FF0000" }}>
        <label htmlFor="CreatorName">Creator Name:</label>
        <input type="text" id="CreatorName" name="CreatorName" defaultValue="Blam" />

        <label htmlFor="GameCode">Game Code:</label>
        <input type="text" id="GameCode" name="GameCode" defaultValue="Blam" />

        <input type="submit" value="Join" />
      </form>

      <form id="startGame" style={{ fontSize: 14, textAlign: "center", backgroundColor: "#4BC445" }}>
        <input type="submit" value="Start Game" />
      </form>

      <form id="makeMove" style={{ fontSize: 14, textAlign: "center", backgroundColor: "#F9D165" }}>
        <label htmlFor="CreatorName">PlayerName:</label>
        <input type="text" id="CreatorName" name="CreatorName" defaultValue="Blam" disabled />

        <label htmlFor="GameCode">Game Code:</label>
        <input type="text" id="GameCode" name="GameCode" defaultValue="Blam" disabled />

        <label htmlFor="betValue">Bet Value:</label>
        <input type="text" id="betValue" name="betValue" defaultValue="Blam" />

        <label htmlFor="betCount">Bet Count:</label>
        <input type="text" id="betCount" name="betCount" defaultValue="Blam" />

        <input type="submit" value="Submit" />
      </form>

      <GameStatus game={game} />
    </div>
  );
}

export default App;
