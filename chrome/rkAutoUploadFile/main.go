package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/767829413/advanced-go/chrome"

	"github.com/chromedp/chromedp"
)

func main() {
	// 创建基础Chrome上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-sync", false),
		chromedp.Flag("headless", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 创建浏览器实例
	// ctx, cancel = chromedp.NewContext(ctx)
	// defer cancel()

	var wg sync.WaitGroup
	numTabs := 10

	for i := 0; i < numTabs; i++ {
		wg.Add(1)
		go func(tabIndex int) {
			defer wg.Done()
			// 为每个标签页创建新的上下文
			tabCtx, cancel := chromedp.NewContext(ctx)
			defer cancel()
			uploadInTab(tabCtx, tabIndex)
		}(i)
	}

	wg.Wait()
	log.Println("所有上传任务完成")
}

func uploadInTab(ctx context.Context, tabIndex int) {
	// 修改日志输出，包含标签页索引
	log.Printf("标签页 %d: 开始上传", tabIndex)

	// 创建一个新的chromedp上下文，配置为非无头模式运行
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-sync", false),
		chromedp.Flag("headless", false), // 设置为false以禁用无头模式,方便调试
	)
	ctx, cancel := context.WithTimeout(ctx, 3600*time.Second)
	defer cancel()
	ctx, cancelChrome := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelChrome()

	ctx, cancelTab := chromedp.NewContext(ctx)
	defer cancelTab()

	// 设置退出机制
	ctx, cancelEnd := context.WithCancel(ctx)
	defer cancelEnd()
	// 定义需要账号登录
	passData, err := os.ReadFile(
		"/home/fangyuan/code/go/src/github.com/767829413/advanced-go/tmp/pass.json",
	)
	if err != nil {
		log.Fatal(err)
	}
	var pass chrome.Pass
	err = json.Unmarshal(passData, &pass)
	if err != nil {
		log.Fatal(err)
	}
	username := pass.Username
	password := pass.Password

	// 定位的记录值
	var dataRowKey string
	// 指定要上传的文件路径
	filePath := filepath.Join("/home/fangyuan/文档/tmp-doc/test_doc/", "画布管理-20241101.xlsx")
	// 结束通知
	endSingal := make(chan struct{})
	go func() {
		// 执行任务
		err = chromedp.Run(
			ctx,
			// 导航到登录页面
			chromedp.Navigate(
				pass.Url,
			),
			chromedp.Sleep(1*time.Second),
			// 等待用户名输入框可见
			chromedp.WaitVisible(`#loginName`, chromedp.ByID),
			// 等待1秒,让内容自动输出
			chromedp.Sleep(1*time.Second),
			// 清空输入框
			chromedp.Clear(`#loginName`, chromedp.ByID),
			// 填写用户名
			chromedp.SendKeys(`#loginName`, username, chromedp.ByID),
			// 添加一个短暂的延迟，以确保输入完成
			chromedp.Sleep(1*time.Second),
			// 等待密码输入框可见
			chromedp.WaitVisible(`#password`, chromedp.ByID),
			// 清空密码输入框
			chromedp.Clear(`#password`, chromedp.ByID),
			// 填写密码
			chromedp.SendKeys(`#password`, password, chromedp.ByID),
			// 添加一个短暂的延迟，以确保输入完成
			chromedp.Sleep(1*time.Second),

			// 等待登录按钮出现
			chromedp.WaitVisible(
				`button.ant-btn.css-ph9edi.ant-btn-primary.ant-btn-lg`,
				chromedp.ByQuery,
			),

			// 尝试移除disabled属性
			chromedp.Evaluate(`
        var btn = document.querySelector('button.ant-btn.css-ph9edi.ant-btn-primary.ant-btn-lg');
        if (btn) {
            btn.disabled = false;
            btn.classList.remove('disabledBtn___xo0Vn');
        }
    `, nil),

			// 短暂等待以确保JS执行完毕
			chromedp.Sleep(1*time.Second),

			// 尝试点击按钮
			chromedp.Click(
				`button.ant-btn.css-ph9edi.ant-btn-primary.ant-btn-lg`,
				chromedp.ByQuery,
			),

			// 等待一段时间，确保点击事件被处理
			chromedp.Sleep(2*time.Second),

			// 如果登录成功，导航到新页面
			chromedp.Navigate(
				pass.Url,
			),

			// 等待新页面加载完成
			chromedp.WaitReady(`body`, chromedp.ByQuery),

			// 等待表格行出现
			chromedp.WaitVisible(`tr[data-row-key]`, chromedp.ByQuery),

			// 短暂等待以确保JS执行完毕
			chromedp.Sleep(1*time.Second),

			// 获取第一个匹配元素的 data-row-key 值
			chromedp.Evaluate(`
			   document.querySelector('tr[data-row-key]').getAttribute('data-row-key')
		   `, &dataRowKey),

			// 等待一小段时间，确保点击事件被处理
			chromedp.Sleep(1*time.Second),
			// 打印跳转的页面
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Printf(
					"获取到的 data-row-key 值: %s,即将跳转的页面 %s",
					dataRowKey,
					pass.Url+pass.TargetSuffix+dataRowKey,
				)
				return nil
			}),
			// 短暂等待以确保JS执行完毕
			chromedp.Sleep(1*time.Second),

			// 导航到指定页面
			// 定义一个函数来设置和验证URL
			chromedp.ActionFunc(func(ctx context.Context) error {
				targetURL := pass.Url + pass.TargetSuffix + dataRowKey
				maxAttempts := 5
				// 检查当前URL
				var currentURL string
				for attempt := 0; attempt < maxAttempts; attempt++ {
					// 设置URL
					if err := chromedp.Evaluate(`window.location.href = '`+targetURL+`';`, nil).Do(ctx); err != nil {
						return err
					}

					// 等待页面加载
					if err := chromedp.WaitReady(`body`, chromedp.ByQuery).Do(ctx); err != nil {
						return err
					}

					// 等待可能的重定向或路由变化
					time.Sleep(2 * time.Second)

					if err := chromedp.Evaluate(`window.location.href`, &currentURL).Do(ctx); err != nil {
						return err
					}

					log.Printf("尝试 %d: 当前URL: %s", attempt+1, currentURL)

					if currentURL == targetURL {
						log.Printf("URL设置成功")
						return nil
					}

					log.Printf("URL不匹配，重试...")
				}

				return fmt.Errorf("无法设置正确的URL，最后的URL: %s", currentURL)
			}),

			// 等待新页面加载完成
			chromedp.WaitReady(`body`, chromedp.ByQuery),

			chromedp.Sleep(1*time.Second),

			// 等待上传按钮出现
			chromedp.WaitVisible(`button.ant-btn.css-ph9edi.ant-btn-primary`, chromedp.ByQuery),

			// 等待一小段时间，确保文件被选择
			chromedp.Sleep(1*time.Second),
			// 直接设置文件路径
			chromedp.SetUploadFiles(
				`input[type="file"]`,
				[]string{filePath},
				chromedp.ByQuery,
			),
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Printf("准备点击上传按钮")
				return nil
			}),
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Printf("点击上传按钮完成")
				return nil
			}),
			chromedp.Sleep(2*time.Second), // 等待一段时间，确保点击事件被处理
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Printf("等待2秒后")
				return nil
			}),
			// 结束
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Printf("开始执行结束逻辑")
				endSingal <- struct{}{}
				return nil
			}),
		)
		if err != nil {
			log.Fatal("chromedp.Run: ", err)
		}
	}()
	<-endSingal
	log.Println("上传成功")

	log.Printf("标签页 %d: 上传完成", tabIndex)
}
