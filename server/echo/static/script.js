// ここはVueそのまま
var app = new Vue({
    el: "#app",
    data: function() {
        return {
            message: 'Hello World!'
        }
    },
    router: new VueRouter({
        // 各ルートにコンポーネントをマッピングします
        // コンポーネントはVue.extend() によって作られたコンポーネントコンストラクタでも
        // コンポーネントオプションのオブジェクトでも構いません
        routes: [
            {
                path: '/top',
                component: {
                    template: '<div>トップページです。</div>'
                }
            },
            {
                path: '/users',
                component: {
                    template: '<div>ユーザー一覧ページです。</div>'
                }
            },
            {
                path: '/todo',
                component: {
                    template: '#todo-item'
                }
            }
        ]
    })
}).$mount('#app');