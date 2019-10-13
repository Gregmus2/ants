import {App} from "./App";
import {Tile} from "./Tile";

App.init();

fetch('http://127.0.0.1:12301/size')
    .then(response => response.text())
    .then(body => {
        let size = parseInt(body);
        let tileSize = Math.round(App.canvas.width / size)
        for (let x = 0; x < size; x++) {
            App.elements[x] = [];
            for (let y = 0; y < size; y++) {
                App.elements[x][y] = new Tile(
                    Math.round(x * tileSize),
                    Math.round(y * tileSize),
                    tileSize,
                    tileSize,
                    "#FFFFFF"
                );
            }
        }
        App.loop();
    });
