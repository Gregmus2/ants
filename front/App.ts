import {Tile} from "./Tile";

export class App{
    static canvas: HTMLCanvasElement;
    static ctx: CanvasRenderingContext2D;
    static elements: Array<Array<Tile>> = [];
    static id: string = null;
    static part: number = 1;
    static select: HTMLSelectElement;
    static buffer: Array<Array<Array<string>>> = [];

    static init(){
        App.createSelect();
        App.updateSelect();
        App.canvas = document.createElement("canvas");
        let size = window.innerWidth < window.innerHeight ? window.innerWidth : window.innerHeight;
        App.canvas.width = size;
        App.canvas.height = size;
        App.canvas.style.border = '1px solid #000';
        App.canvas.addEventListener('contextmenu', event => event.preventDefault());
        document.body.appendChild(App.canvas);
        App.ctx = App.canvas.getContext('2d');
    }
    
    static loop(){
        if (App.id == null && App.buffer.length == 0) {
            return;
        }

        App.update();

        for (let row of App.elements) {
            for (let el of row) {
                if (el.changed) {
                    App.ctx.clearRect(el.x, el.y, el.w, el.h);
                    App.ctx.drawImage(el.draw(), el.x, el.y);
                }
            }
        }

        requestAnimationFrame(App.loop);
    }
    
    static update() {
        if (App.buffer.length < 200 && App.id != null) {
            App.fetchMap();
        }

        App.renderFromBuffer();
    }

    // todo что делать с порядком вставки в буфер
    static fetchMap() {
        let prom = fetch('/get?id=' + App.id + '&part=' + App.part)
            .then(response => {
                switch (response.status) {
                    case 200:
                        return response.json();
                    case 404:
                        return Promise.reject('404')
                }
            })
            .then(body => {
                App.buffer = App.buffer.concat(body);
            }).catch(() => {App.id = null});
        App.part++;

        return prom;
    }

    static createSelect(){
        App.select = document.createElement("select");
        document.body.appendChild(App.select);
        App.select.addEventListener("change", function() {
            App.id = this.selectedOptions.item(0).value;
            App.part = 1;
            App.fetchMap().then(() => App.loop());
        })
    }

    static updateSelect(){
        if (App.select == undefined) {
            return;
        }

        fetch('/pipes')
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
}