module.exports = {
    content: ["./templates/**/*.{tmpl,js}" ,"./public/**/*.{tmpl,js}"],
    theme: {
        extend: {
            keyframes: {
                'fade-in-move': {
                    '0%': { opacity: '0', transform: 'translateY(10px)' },
                    '100%': { opacity: '1', transform: 'translateY(0)' },
                },
            },
            animation: {
                'fade-in-move': 'fade-in-move 0.3s ease-out forwards',
            },
        },
    },
    plugins: [],
}
