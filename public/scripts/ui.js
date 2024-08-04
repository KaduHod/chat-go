const menuHamburguer = document.getElementById('menu-hamburguer')
const salaAberta = document.getElementById("chat-salaaberta")
const menuLateral = document.getElementById("chat-salas")
console.log({menuHamburguer, salaAberta, menuLateral})
menuHamburguer.addEventListener('click', () => {
    menuLateral.classList.toggle("hidden")
    salaAberta.classList.toggle("hidden")
    window.scrollTo({top: document.body.scrollHeight, behavior: "smooth"})
})
