export type DrinkMetadata = {
	title: string;
	category: string;
	slug: string;
	tags?: string[];
	date?: string;
};

export type DrinkPage = {
	metadata: DrinkMetadata;
	content: string;
};
