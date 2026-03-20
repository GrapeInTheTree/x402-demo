package practice

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
)

// EIP3009FlowModel manages the EIP-3009 10-step payment flow.
type EIP3009FlowModel struct {
	sm     stepManager
	width  int
	height int
	cfg    *config.ExplorerConfig
}

func NewEIP3009FlowModel(width, height int, cfg *config.ExplorerConfig) *EIP3009FlowModel {
	descriptions := []stepDesc{
		{"지갑 주소 확인", "—", "—"},
		{"—", "GET /supported 호출", "/supported 응답"},
		{"GET /weather (결제 없음)", "402 반환", "—"},
		{"PAYMENT-REQUIRED 디코딩", "—", "—"},
		{"EIP-712 서명 생성", "—", "—"},
		{"PAYMENT-SIGNATURE 전송", "헤더 파싱 → /verify", "—"},
		{"—", "/verify 요청 전달", "서명/잔액/시뮬레이션 검증"},
		{"200 + 데이터 수신", "데이터 반환 + /settle", "—"},
		{"—", "—", "온체인 트랜잭션 제출"},
		{"최종 잔액 확인", "—", "—"},
	}

	return &EIP3009FlowModel{
		sm:     newStepManager(demo.NewFlowState("eip3009"), descriptions),
		width:  width,
		height: height,
		cfg:    cfg,
	}
}

func (m *EIP3009FlowModel) Update(msg tea.Msg) tea.Cmd {
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

func (m *EIP3009FlowModel) View() string {
	return m.sm.view(m.width)
}
