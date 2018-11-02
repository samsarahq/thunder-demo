import React from 'react';
import { connectGraphQL, mutate } from "thunder-react/lib/store";

class NewGame extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            wsConnected: false
        }
    }

    createGame = async () => {
        try {
        const game = await mutate({
            query: `{createGame(name: $name) {id}}`,
            variables: {
                name: "It's Thundoku Time!"
            }
        });
        window.location.pathname = `/${game.createGame.id}`;
    }
    catch (e) {
        console.log(e);
    }
    }

    render = () => {
        if (!this.state.wsConnected) {
            this.setState({ wsConnected: true });
            this.createGame();
        }
        return <div/>
    }
}


export default connectGraphQL(NewGame, () => ({
    query: "{games {id}}",
    variables: {},
    onlyValidData: true
}));