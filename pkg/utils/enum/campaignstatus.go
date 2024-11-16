package enum

type CampaignStatus struct {
	Status string
	Color  string
}

func StatusDelevering() CampaignStatus {
	return CampaignStatus{
		Status: "Delivering",
		Color:  "#26B50F",
	}
}

func StatusPaused() CampaignStatus {
	return CampaignStatus{
		Status: "Paused",
		Color:  "#5E5E5E",
	}
}

func StatusOutOfBudget() CampaignStatus {
	return CampaignStatus{
		Status: "Out of budget",
		Color:  "#9D4505",
	}
}

func StatusArchived() CampaignStatus {
	return CampaignStatus{
		Status: "Archived",
		Color:  "#C2C2C2",
	}
}
