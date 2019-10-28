import {Tile} from "./Tile";

export class App{
    static canvas: HTMLCanvasElement;
    static ctx: CanvasRenderingContext2D;
    static elements: Array<Array<Tile>> = [];
    static id: string = null;
    static part: number = 1;
    static select: HTMLSelectElement;
    static buffer: Array<Array<Array<string>>> = [];
    static updating: boolean = false;
    static playersContainer: HTMLElement;

    static init(){
        App.createSelect();
        App.updateSelect();
        App.initForm();
        App.createPlayers();
        App.updatePlayers();
        App.canvas = document.querySelector("canvas");
        let size = window.innerWidth < window.innerHeight ? window.innerWidth : window.innerHeight;
        App.canvas.width = size;
        App.canvas.height = size;
        App.canvas.style.borderRight = '1px solid #000';
        App.canvas.style.borderLeft= '1px solid #000';
        App.canvas.addEventListener('contextmenu', event => event.preventDefault());
        document.body.appendChild(App.canvas);
        App.ctx = App.canvas.getContext('2d');
        App.ctx.imageSmoothingEnabled = false;
    }
    
    static loop(){
        if (App.id == null && App.buffer.length == 0) {
            return;
        }

        App.update();

        for (let row of App.elements) {
            for (let el of row) {
                if (el.changed) {
                    el.draw()
                }
            }
        }

        requestAnimationFrame(App.loop);
        // setTimeout(() => App.loop(), 500)
    }
    
    static update() {
        if (App.buffer.length < 200 && App.id != null && !App.updating) {
            App.fetchMap();
        }

        App.renderFromBuffer();
    }

    static fetchMap() {
        App.updating = true;
        let prom = fetch('/api/get?id=' + App.id + '&part=' + App.part)
            .then(response => {
                switch (response.status) {
                    case 200:
                        return response.json();
                    case 404:
                        return Promise.reject('404')
                }
            })
            .then(body => {
                App.updating = false;
                if (body.length == 0) {
                    App.id = null;
                    return;
                }

                App.buffer = App.buffer.concat(body);
            }).catch(() => {App.id = null});
        App.part++;

        return prom;
    }

    static createSelect(){
        App.select = document.querySelector("select");
        App.select.addEventListener("change", function() {
            // @ts-ignore
            UIkit.offcanvas(document.getElementById('game')).hide();
            App.buffer = [];
            App.id = this.selectedOptions.item(0).value;
            App.part = 1;
            App.fetchMap().then(() => App.loop());
        })
    }

    static updateSelect(){
        if (App.select == undefined) {
            return;
        }

        fetch('/api/pipes')
            .then(response => response.json())
            .then(body => {
                for (let index in App.select.options) {
                    App.select.remove(Number(index))
                }

                let opt = document.createElement("option");
                opt.value = undefined;
                opt.text = '';
                App.select.add(opt, null);

                for (let id of body) {
                    let opt = document.createElement("option");
                    opt.value = id;
                    opt.text = id;
                    App.select.add(opt, null)
                }

                setTimeout(() => (App.updateSelect()), 5000);
            })
    }

    static createPlayers(){
        App.playersContainer = document.getElementById('players');
        document.getElementById('refresh').addEventListener("click", () => {
            this.updatePlayers();
        });
        document.getElementById('start').addEventListener("click", () => {
            this.startGame();
        });
    }

    static startGame() {
        const request = new XMLHttpRequest();
        request.open('POST', '/api/start', true);
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                App.updateSelect();
                alert("Your number: " + this.response)
            } else {
                alert(this.response)
            }
        };
        request.onerror = function (err) {
            alert(err)
        };

        let names = [];
        App.playersContainer.querySelectorAll(':checked').forEach(function (value) {
            names.push(value.nextSibling.textContent)
        });

        request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
        request.send('names=' + names.join(','));
        App.updatePlayers()
    }

    static updatePlayers() {
        if (App.playersContainer == undefined) {
            return;
        }

        fetch('/api/players')
            .then(response => response.json())
            .then(body => {
                App.playersContainer.innerHTML = '';

                for (let name of body) {
                    let div = document.createElement("div");
                    div.className = 'uk-margin uk-grid-small uk-child-width-auto uk-grid';
                    let label = document.createElement("label");
                    let span = document.createElement("span");
                    span.innerText = name;
                    span.className = 'uk-margin-small-left';
                    let input = document.createElement("input");
                    input.className = 'uk-checkbox';
                    input.type = 'checkbox';

                    div.appendChild(label);
                    label.appendChild(input);
                    label.appendChild(span);
                    App.playersContainer.appendChild(div)
                }
            })
    }

    static renderFromBuffer(){
        let map = App.buffer.shift();
        if (!map) {
            return;
        }

        App.renderMap(map)
    }


    static renderMap(map: Array<Array<string>>){
        for (let x in map) {
            for (let y in map) {
                App.elements[x][y].setColor(map[x][y] == "" ? 'black' : map[x][y])
            }
        }
    }

    static initForm(){
        let form = document.querySelector('form');
        form.addEventListener('submit', (event) => {event.preventDefault(); this.submitForm(form)});
    }

    static submitForm(form: HTMLFormElement){
        const request = new XMLHttpRequest();
        request.open('POST', '/api/register', true);
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                form.reset();
                App.updatePlayers();
                // @ts-ignore
                UIkit.offcanvas(document.getElementById('register')).hide();
                // @ts-ignore
                UIkit.offcanvas(document.getElementById('game')).show();
            } else {
                alert(this.response)
            }
        };
        request.onerror = function (err) {
            alert(err)
        };
        request.send(new FormData(form));
    }
}