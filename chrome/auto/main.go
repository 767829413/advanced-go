package main

import (
	"context"
	"github.com/767829413/advanced-go/chrome"
	"log"
	"time"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

// <div style="top:4px;height:24px;" class="view-line"><span><span></span></span></div>
func main() {
	AutoSearchAliyunLog(context.Background())
}

func AutoSearchAliyunLog(c context.Context) {
	// 创建一个新的chromedp上下文，配置为非无头模式运行
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.UserDataDir("C:/Users/NUC/AppData/Local/Google/Chrome/User Data"),
		chromedp.Flag("disable-sync", false),
		chromedp.Flag("headless", false), // 设置为false以禁用无头模式,方便调试
	)
	ctx, cancel := chromedp.NewExecAllocator(c, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	// 定义阿里云ram用户名和密码
	username := ""
	password := ""

	// 生成2fa-totp code
	// 使用自己的secret
	code, err := chrome.GetTotpCode(
		"",
	)
	if err != nil {
		log.Fatal(err)
	}
	// 执行任务
	err = chromedp.Run(
		ctx,
		// 导航到登录页面
		chromedp.Navigate(
			`https://signin.aliyun.com/plaso_rk.onaliyun.com/login.htm?callback=https%3A%2F%2Fsls.console.aliyun.com%2Flognext%2Fproject%2Frongke-slb%2Flogsearch%2Frongke-slb-log%3FslsRegion%3Dcn-hangzhou&accounttraceid=c9b4935409a4441f801572d6c0c46aabrbke&cspNonce=QvU41wVAK6&oauth_callback=https%3A%2F%2Fsls.console.aliyun.com%2Flognext%2Fproject%2Frongke-slb%2Flogsearch%2Frongke-slb-log%3FslsRegion%3Dcn-hangzhou&spma=a2c44&spmb=11131515#/main`,
		),
		chromedp.Sleep(3*time.Second),
		// 等待用户名输入框可见
		chromedp.WaitVisible(`#loginName`, chromedp.ByID),
		// 等待2秒,让内容自动输出
		chromedp.Sleep(1*time.Second),
		// 使用SendKeys发送空格和退格键来清空输入框
		chromedp.Click(`.next-icon.next-icon-close.next-xs`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		// 填写用户名
		chromedp.SendKeys(`#loginName`, username, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		// 滑动
		// chromedp.ActionFunc(SwipeRight),
		// chromedp.Sleep(1*time.Second),
		// 点击下一步按钮
		chromedp.Click(`//button/span[contains(text(), '下一步')]`, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.WaitVisible(`#loginPassword`, chromedp.ByID),
		chromedp.SendKeys(`#loginPassword`, password, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		// 点击登录按钮
		chromedp.Click(`button.next-btn.next-large.next-btn-primary`, chromedp.ByQuery),
		// 输入二次验证
		chromedp.SendKeys(`input[placeholder="请输入 6 位数字安全码"]`, code, chromedp.ByQuery),
		// 等待一些时间以确保登录成功，这里等待的时间可能需要根据实际情况调整
		// 点击提交验证
		chromedp.Click(`button.next-btn.next-large.next-btn-primary`, chromedp.ByQuery),
		// 确保目标元素可见
		chromedp.WaitVisible(`div.view-line > span > span`, chromedp.ByQuery),
		// 向目标元素发送文字
		// 使用JavaScript替换最里面的<span></span>为<span class="mtk1">error</span>
		chromedp.Evaluate(
			`document.querySelector('div.view-line > span > span').innerHTML = '<span class="mtk1">error</span>';`,
			nil,
		),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`div.BaseSearch-m__search-btn__525b93f4`, chromedp.ByQuery),
		chromedp.Sleep(300*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("登录成功")
}

// 模拟向右滑动
func SwipeRight(ctx context.Context) error {
	var x, y, width, height float64
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`(() => {
        const rect = document.querySelector('#nc_1_n1z').getBoundingClientRect();
        return [rect.x, rect.y, rect.width, rect.height];
    })()`, &[]interface{}{&x, &y, &width, &height}),
	)
	if err != nil {
		log.Fatal(err)
	}

	startX, startY := x, y // 假设的起始点坐标
	endX := startX + 200   // 假设向右滑动200像素为目标位置

	// 模拟鼠标按下动作
	if err := input.DispatchMouseEvent(input.MousePressed, float64(startX), float64(startY)).Do(ctx); err != nil {
		return err
	}

	// 模拟鼠标滑动过程，可以通过增加中间步骤来模拟更自然的手动滑动
	steps := 10 // 定义滑动步骤数
	for i := 1; i <= steps; i++ {
		x := startX + (endX-startX)*float64(i)/float64(steps)
		if err := input.DispatchMouseEvent(input.MouseMoved, float64(x), float64(startY)).Do(ctx); err != nil {
			return err
		}
		// 增加小的延迟来模拟人手滑动的停顿
		time.Sleep(50 * time.Millisecond) // 50毫秒的停顿
	}

	// 模拟鼠标释放动作
	if err := input.DispatchMouseEvent(input.MouseReleased, float64(endX), float64(startY)).Do(ctx); err != nil {
		return err
	}

	return nil
}
