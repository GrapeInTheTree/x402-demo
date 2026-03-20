package practice

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
)

// Permit2FlowModel manages the Permit2 10-step payment flow.
type Permit2FlowModel struct {
	sm     stepManager
	width  int
	height int
	cfg    *config.ExplorerConfig
}

func NewPermit2FlowModel(width, height int, cfg *config.ExplorerConfig) *Permit2FlowModel {
	descriptions := []stepDesc{
		{"지갑 주소 + Permit2 approve 확인", "—", "—"},
		{"—", "GET /supported 호출", "/supported 응답"},
		{"GET /weather (결제 없음)", "402 + assetTransferMethod:permit2", "—"},
		{"PAYMENT-REQUIRED 디코딩 (Permit2)", "—", "—"},
		{"Permit2 EIP-712 서명 생성", "—", "—"},
		{"PAYMENT-SIGNATURE 전송", "헤더 파싱 → /verify", "—"},
		{"—", "/verify 요청 전달", "Permit2 서명 + allowance 검증"},
		{"200 + 데이터 수신", "데이터 반환 + /settle", "—"},
		{"—", "—", "x402Permit2Proxy.settle() 제출"},
		{"최종 잔액 확인", "—", "—"},
	}

	return &Permit2FlowModel{
		sm:     newStepManager(demo.NewFlowState("permit2"), descriptions),
		width:  width,
		height: height,
		cfg:    cfg,
	}
}

func (m *Permit2FlowModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "n":
			m.sm.next()
		case "p":
			m.sm.prev()
		}
	}
	return nil
}

func (m *Permit2FlowModel) View() string {
	return m.sm.view(m.width)
}
