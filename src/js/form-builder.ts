import { Formio, Templates } from '@formio/js';
import { validation } from './validation';

Templates.framework = 'goforms'; // Ensure this is actually required

interface FormBuilderOptions {
  disabled?: string[];
  noNewEdit?: boolean;
  noDefaultSubmitButton?: boolean;
  alwaysConfirmComponentRemoval?: boolean;
  formConfig?: Record<string, any>;
  resourceTag?: string;
  editForm?: Record<string, any>;
  language?: string;
  builder?: object;
  display?: 'form' | 'wizard' | 'pdf';
  resourceFilter?: string;
  noSource?: boolean;
  showFullJsonSchema?: boolean;
}

export class FormBuilder {
  private container: HTMLElement;
  private builder!: Formio.BuilderInstance;
  private formId: number;
  private currentSchema: Record<string, any> = { display: 'form', components: [] };

  constructor(containerId: string, formId: number) {
    console.log(`Initializing FormBuilder with formId: ${formId}`);
    this.container = this.getContainer(containerId);
    this.formId = this.validateFormId(formId);
    this.init();
  }

  private getContainer(containerId: string): HTMLElement {
    const container = document.getElementById(containerId);
    if (!container) throw new Error(`Error: Container '${containerId}' not found.`);
    return container;
  }

  private validateFormId(formId: number): number {
    if (isNaN(formId) || formId < 0) throw new Error(`Error: Invalid form ID '${formId}'.`);
    return formId;
  }

  private async init() {
    const builderOptions: FormBuilderOptions = {
      display: 'form',
      noDefaultSubmitButton: true,
      builder: {
        basic: { title: 'Basic Fields', default: true, weight: 0, components: { textfield: true } },
      }
    };

    try {
      const schema = await this.loadExistingSchema();
      this.builder = await Formio.builder(this.container, schema, builderOptions);
    } catch (error) {
      console.error('Initialization failed:', error);
    }
  }

  private async loadExistingSchema(): Promise<any> {
    if (this.formId === 0) return { display: 'form', components: [] };

    try {
      const response = await validation.fetchWithCSRF(`/dashboard/forms/${this.formId}/schema`, { method: 'GET' });
      if (!response.ok) throw new Error(`Failed to load schema: ${response.status} ${response.statusText}`);

      const schema = await response.json();
      console.log(`Loaded form schema for form ID: ${this.formId}`, schema);
      return schema; // ✅ Return the schema properly
    } catch (error) {
      console.error('Error loading form schema:', error);
      return { display: 'form', components: [] }; // Fallback schema
    }
  }

  public async saveSchema(): Promise<boolean> {
    try {
      const formioSchema = this.builder.schema;
      const response = await validation.fetchWithCSRF(`/dashboard/forms/${this.formId}/schema`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formioSchema)
      });

      if (!response.ok) throw new Error('Failed to save schema');

      console.log('Schema saved successfully:', formioSchema);
      this.currentSchema = formioSchema;
      return true;
    } catch (error) {
      console.error('Error saving form schema:', error);
      return false;
    }
  }
}