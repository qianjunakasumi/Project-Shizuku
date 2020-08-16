const dom = new Vue({
    el: '#wrapper',
    data: {
        followers: 0,
        new_followers: '',
        huanbi_yesterday: '',
        huanbi_rate_yesterday: '',
        dingji_yuechu: '',
        tongbi_dingji_rate_shangyue: '',
    },
})

function sign(num) {

    if (num >= 0) {

        return '+' + num

    }

    return num

}

function chart(data) {

    const chart = new G2.Chart({
        container: 'card-chart',
        autoFit: true,
        height: 180,
    })

    chart.animate(false)
    chart.data(data)
    chart.scale({
        value: {
            nice: true,
        },
        date: {
            range: [0, 1],
        },
    })
    chart.axis('value', {
        label: {
            formatter: (value) => {

                return value / 1000 + 'k';

            },
        },
    })
    chart
        .point()
        .position('date*value')
        .label('date*value', (date, value) => {
            return {
                content: value,
            }
        })

    chart.area().position('date*value').color('#ff5760')
    chart.line().position('date*value').color('#FF5760')
    chart.render()

}

function init() {

    dom.followers = data.followersCount

    dom.new_followers = sign(data.newFollowersCount)

    dom.huanbi_yesterday = sign(data.huanbiRate)

    dom.huanbi_rate_yesterday = sign(data.ydayHuanbiRate)

    dom.dingji_yuechu = sign(data.dingjiRate)

    dom.tongbi_dingji_rate_shangyue = sign(data.syueDingjiRate)

    chart(data.per)

}

let data

fetch('data.json')
    .then(res => {

        return res.json()

    })
    .then(d => {

        data = d
        console.log(d)

        init()

    })
