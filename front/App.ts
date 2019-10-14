import {Tile} from "./Tile";

export class App{
    static canvas: HTMLCanvasElement;
    static ctx: CanvasRenderingContext2D;
    static elements: Array<Array<Tile>> = [];
    static pipe: number;
    static select: HTMLSelectElement;

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
        App.update();
        App.updateSelect();

        for (let row of App.elements) {
            for (let el of row) {
                if (el.changed) {
                    App.ctx.clearRect(el.x, el.y, el.w, el.h);
                    App.ctx.drawImage(el.draw(), el.x, el.y);
                }
            }
        }
    }
    
    static update(){
        if (App.pipe == undefined) {
            setTimeout(() => (App.loop()), 10000);
            return;
        }

        fetch('/get?id=' + App.pipe)
            .then(response => {
                switch (response.status) {
                    case 200:
                        return response.json();
                    case 404:
                        return Promise.reject('404')
                }
            })
            .then(body => {
                for (let x in body) {
                    for (let y in body) {
                        App.elements[x][y].setColor(body[x][y])
                    }
                }
    
                setTimeout(() => (App.loop()), 50);
            })
            .catch(() => {setTimeout(() => (App.loop()), 10000); console.log(123)});
    }

    static createSelect(){
        App.select = document.createElement("select");
        document.body.appendChild(App.select);
        App.select.addEventListener('onchange', function() {
            App.pipe = this.selectedOptions.item(0).value;
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
                for (let id of body) {
                    let opt = document.createElement("option");
                    opt.value = id;
                    opt.text = id;
                    App.select.add(opt, null)
                }
            })
    }
}